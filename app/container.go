package container

import (
	"clean/app/http/controllers"
	"clean/app/http/middlewares"
	repoImpl "clean/app/repository/impl"
	svcImpl "clean/app/svc/impl"
	"clean/infra/conn"
)

func Init(g interface{}) {
	db := conn.Db()
	redis := conn.Redis()
	acl := middlewares.ACL

	// register all repos impl, services impl, controllers
	sysRepo := repoImpl.NewSystemRepository(db, redis)
	userRepo := repoImpl.NewMySqlUsersRepository(db)
	roleRepo := repoImpl.NewMySqlRolesRepository(db)
	permissionRepo := repoImpl.NewMySqlPermissionsRepository(db)
	speRepo := repoImpl.NewSpecializationRepository(db)

	cacheSvc := svcImpl.NewRedisService(redis)
	sysSvc := svcImpl.NewSystemService(sysRepo)
	userSvc := svcImpl.NewUsersService(userRepo, cacheSvc)
	tokenSvc := svcImpl.NewTokenService(userRepo, cacheSvc)
	authSvc := svcImpl.NewAuthService(userRepo, cacheSvc, tokenSvc)
	roleSvc := svcImpl.NewRolesService(roleRepo)
	permissionSvc := svcImpl.NewPermissionsService(permissionRepo)
	spSvc := svcImpl.NewSpecializationService(speRepo)

	controllers.NewSystemController(g, sysSvc)
	controllers.NewAuthController(g, authSvc, userSvc)
	controllers.NewUsersController(g, userSvc)
	controllers.NewRolesController(g, acl, roleSvc)
	controllers.NewPermissionsController(g, acl, permissionSvc)
	controllers.NewSpecializationController(g, spSvc)
}
