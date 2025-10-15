// Ejemplo de consumo del servicio SCP Go Identity desde aplicación .NET en IIS
// Este código va en tu aplicación .NET que corre en IIS

using System;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;

public class SCPIdentityService
{
    private readonly HttpClient _httpClient;
    private readonly string _serviceBaseUrl;

    public SCPIdentityService(string serviceBaseUrl = "http://localhost:8081")
    {
        _serviceBaseUrl = serviceBaseUrl;
        _httpClient = new HttpClient()
        {
            Timeout = TimeSpan.FromSeconds(30)
        };
    }

    public async Task<AuthenticationResult> AuthenticateAsync(string email, string password, bool generateToken = true)
    {
        var request = new
        {
            email = email,
            password = password,
            generateToken = generateToken
        };

        try
        {
            var json = JsonConvert.SerializeObject(request);
            var content = new StringContent(json, Encoding.UTF8, "application/json");

            var response = await _httpClient.PostAsync($"{_serviceBaseUrl}/authenticate", content);
            var responseContent = await response.Content.ReadAsStringAsync();

            if (response.IsSuccessStatusCode)
            {
                return JsonConvert.DeserializeObject<AuthenticationResult>(responseContent);
            }
            else
            {
                var errorResult = JsonConvert.DeserializeObject<AuthenticationResult>(responseContent);
                return new AuthenticationResult 
                { 
                    Success = false, 
                    Message = errorResult?.Message ?? $"HTTP {response.StatusCode}" 
                };
            }
        }
        catch (HttpRequestException ex)
        {
            return new AuthenticationResult 
            { 
                Success = false, 
                Message = $"Error de conexión al servicio de identidad: {ex.Message}" 
            };
        }
        catch (TaskCanceledException ex)
        {
            return new AuthenticationResult 
            { 
                Success = false, 
                Message = "Timeout al conectar con el servicio de identidad" 
            };
        }
        catch (Exception ex)
        {
            return new AuthenticationResult 
            { 
                Success = false, 
                Message = $"Error inesperado: {ex.Message}" 
            };
        }
    }

    public async Task<bool> IsServiceHealthyAsync()
    {
        try
        {
            var response = await _httpClient.GetAsync($"{_serviceBaseUrl}/health");
            return response.IsSuccessStatusCode;
        }
        catch
        {
            return false;
        }
    }

    public void Dispose()
    {
        _httpClient?.Dispose();
    }
}

// Modelos de respuesta
public class AuthenticationResult
{
    public bool Success { get; set; }
    public string Message { get; set; }
    public UserInfo User { get; set; }
    public string[] Roles { get; set; }
    public ClaimInfo[] Claims { get; set; }
    public string Token { get; set; }
    public DateTime? ExpiresAt { get; set; }
}

public class UserInfo
{
    public string Id { get; set; }
    public string Email { get; set; }
    public string UserName { get; set; }
    public string FirstName { get; set; }
    public string LastName { get; set; }
    public string Culture { get; set; }
    public bool IsEmailConfirmed { get; set; }
    public string PhoneNumber { get; set; }
}

public class ClaimInfo
{
    public string Type { get; set; }
    public string Value { get; set; }
}

// Ejemplo de uso en un controlador MVC
public class AccountController : Controller
{
    private readonly SCPIdentityService _identityService;

    public AccountController()
    {
        _identityService = new SCPIdentityService();
    }

    [HttpPost]
    public async Task<IActionResult> Login(string email, string password)
    {
        if (string.IsNullOrEmpty(email) || string.IsNullOrEmpty(password))
        {
            ViewBag.Error = "Email y contraseña son requeridos";
            return View();
        }

        var result = await _identityService.AuthenticateAsync(email, password, true);

        if (result.Success)
        {
            // Autenticación exitosa
            // Aquí puedes crear la sesión, cookies, etc.
            
            // Ejemplo: guardar info en session
            Session["UserId"] = result.User.Id;
            Session["UserEmail"] = result.User.Email;
            Session["UserRoles"] = result.Roles;
            Session["AuthToken"] = result.Token;

            // Redirigir al dashboard o página principal
            return RedirectToAction("Dashboard", "Home");
        }
        else
        {
            ViewBag.Error = result.Message;
            return View();
        }
    }

    protected override void Dispose(bool disposing)
    {
        if (disposing)
        {
            _identityService?.Dispose();
        }
        base.Dispose(disposing);
    }
}