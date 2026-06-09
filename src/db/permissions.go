package db

type Permission int64

const (
	VIEW_SERVER_STATUS     Permission = 0x1
	START_SERVER           Permission = 0x2
	STOP_SERVER            Permission = 0x4
	RESTART_SERVER         Permission = 0x8
	VIEW_PLAYERS           Permission = 0x10
	VIEW_SERVER_PROPERTIES Permission = 0x20
	READ_SERVER_LOGS       Permission = 0x40
	SEND_COMMANDS          Permission = 0x80
	VIEW_LOGS_HISTORY      Permission = 0x100
	DELETE_LOG             Permission = 0x200
	CREATE_USER            Permission = 0x400
	DELETE_USER            Permission = 0x800
	MODIFY_USER            Permission = 0x1000
	VIEW_USERS             Permission = 0x2000
)

func CheckPermissions(permissions int64, permissionToCheck int64) bool {
	return (permissions & permissionToCheck) == permissionToCheck
}
