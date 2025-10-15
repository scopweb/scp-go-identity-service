-- Script para crear un usuario de prueba en base de datos .NET Identity
-- EJECUTAR SOLO SI NECESITAS UN USUARIO DE PRUEBA

-- 1. Insertar usuario (el password hash es para "Test123!")
INSERT INTO AspNetUsers (
    Id, 
    UserName, 
    NormalizedUserName, 
    Email, 
    NormalizedEmail, 
    EmailConfirmed, 
    PasswordHash, 
    SecurityStamp, 
    ConcurrencyStamp,
    PhoneNumberConfirmed, 
    TwoFactorEnabled, 
    LockoutEnabled, 
    AccessFailedCount,
    FirstName,
    LastName
) VALUES (
    NEWID(),
    'test@scp.com',
    'TEST@SCP.COM',
    'test@scp.com',
    'TEST@SCP.COM',
    1,
    'AQAAAAEAACcQAAAAEJ3/XKJ+z6Z5k5L1q+JxKpA+V7QjLy5G8OmDfJ7mOq9XoA6s2P/4gTZ8xYdV4K1+Uw==', -- Hash para "Test123!"
    NEWID(),
    NEWID(),
    0,
    0,
    1,
    0,
    'Usuario',
    'Prueba'
);

-- 2. Verificar que se creo
SELECT * FROM AspNetUsers WHERE Email = 'test@scp.com';