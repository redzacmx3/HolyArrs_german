package accounts

import "github.com/HolyArrs/holyarrs/app/backend/internal/core/webserver"

// Fork routes (Upkeep, Pick wishes) live here, so scopes.go stays as
// upstream has it and under the file size cap.
func init() {
	for route, perm := range map[string]webserver.Perm{
		"GET /api/v1/upkeep":                          webserver.PermManageLibrary,
		"PUT /api/v1/settings/upkeep":                 webserver.PermSettings,
		"POST /api/v1/upkeep/apply":                   webserver.PermManageLibrary,
		"GET /api/v1/titles/{id}/wishes":              webserver.PermManageLibrary,
		"PUT /api/v1/titles/{id}/wishes":              webserver.PermManageLibrary,
		"POST /api/v1/titles/{id}/wishes/preview":     webserver.PermManageLibrary,
		"GET /api/v1/settings/wishes/{kind}":          webserver.PermSettings,
		"PUT /api/v1/settings/wishes/{kind}":          webserver.PermSettings,
		"POST /api/v1/settings/wishes/{kind}/preview": webserver.PermSettings,
	} {
		RoutePerms[route] = perm
	}
}
