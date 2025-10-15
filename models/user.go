package models

import (
	"database/sql"
)

// AspNetUser represents the AspNetUsers table from .NET Identity
type AspNetUser struct {
	Id                   string         `db:"Id" json:"id"`
	UserName             sql.NullString `db:"UserName" json:"userName"`
	NormalizedUserName   sql.NullString `db:"NormalizedUserName" json:"normalizedUserName"`
	Email                sql.NullString `db:"Email" json:"email"`
	NormalizedEmail      sql.NullString `db:"NormalizedEmail" json:"normalizedEmail"`
	EmailConfirmed       bool           `db:"EmailConfirmed" json:"emailConfirmed"`
	PasswordHash         sql.NullString `db:"PasswordHash" json:"-"`
	SecurityStamp        sql.NullString `db:"SecurityStamp" json:"-"`
	ConcurrencyStamp     sql.NullString `db:"ConcurrencyStamp" json:"-"`
	PhoneNumber          sql.NullString `db:"PhoneNumber" json:"phoneNumber"`
	PhoneNumberConfirmed bool           `db:"PhoneNumberConfirmed" json:"phoneNumberConfirmed"`
	TwoFactorEnabled     bool           `db:"TwoFactorEnabled" json:"twoFactorEnabled"`
	LockoutEnd           sql.NullTime   `db:"LockoutEnd" json:"lockoutEnd"`
	LockoutEnabled       bool           `db:"LockoutEnabled" json:"lockoutEnabled"`
	AccessFailedCount    int            `db:"AccessFailedCount" json:"accessFailedCount"`
	FirstName            sql.NullString `db:"FirstName" json:"firstName"`
	LastName             sql.NullString `db:"LastName" json:"lastName"`
	Culture              sql.NullString `db:"Culture" json:"culture"`
}

// AspNetRole represents the AspNetRoles table from .NET Identity
type AspNetRole struct {
	Id               string         `db:"Id" json:"id"`
	Name             sql.NullString `db:"Name" json:"name"`
	NormalizedName   sql.NullString `db:"NormalizedName" json:"normalizedName"`
	ConcurrencyStamp sql.NullString `db:"ConcurrencyStamp" json:"-"`
}

// AspNetUserRole represents the AspNetUserRoles table from .NET Identity
type AspNetUserRole struct {
	UserId string `db:"UserId" json:"userId"`
	RoleId string `db:"RoleId" json:"roleId"`
}

// AspNetUserClaim represents the AspNetUserClaims table from .NET Identity
type AspNetUserClaim struct {
	Id         int            `db:"Id" json:"id"`
	UserId     string         `db:"UserId" json:"userId"`
	ClaimType  sql.NullString `db:"ClaimType" json:"claimType"`
	ClaimValue sql.NullString `db:"ClaimValue" json:"claimValue"`
}
