package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/utils/methodsutil"
	"clean/app/utils/msgutil"
	"clean/infra/errors"
	"clean/infra/logger"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type users struct {
	*gorm.DB
}

// NewMySqlUsersRepository will create an object that represent the User.Repository implementations
func NewMySqlUsersRepository(db *gorm.DB) repository.IUsers {
	return &users{
		DB: db,
	}
}

func (r *users) Save(user *domain.User) (*domain.User, *errors.RestErr) {
	res := r.DB.Model(&domain.User{}).Create(&user)

	if res.Error != nil {
		logger.Error("error occurred when create user", res.Error)
		restErr := fmt.Sprintf("%s", res.Error)
		return nil, errors.NewInternalServerError(fmt.Sprintf("%s", errors.NewError(restErr)))
	}

	return user, nil
}

func (r *users) GetUser(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr) {
	//var resp domain.User
	var intUser domain.IntermediateUserWithPermissions
	var userWithParams domain.UserWithPerms

	sections := `
		users.*,
		roles.name role_name
	`
	if withPermission {
		sections += ",GROUP_CONCAT(DISTINCT permissions.name) AS permissions"
	}

	query := r.DB.Model(&domain.User{}).
		Select(sections).
		Joins("LEFT JOIN roles ON users.role_id = roles.id").
		Where("users.deleted_at IS NULL")

	if withPermission {
		query = query.
			Joins("JOIN role_permissions ON users.role_id = role_permissions.role_id").
			Joins("JOIN permissions ON role_permissions.permission_id = permissions.id")
	}

	query.Group("users.id")

	res := query.Where("users.id = ?", userID).Find(&intUser)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("getting user with permission"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	err := methodsutil.StructToStruct(intUser, &userWithParams.User)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	userWithParams.RoleName = intUser.RoleName

	if withPermission {
		userWithParams.Permissions = strings.Split(intUser.Permissions, ",")
	}

	return &userWithParams, nil
}

func (r *users) GetUserByID(userID uint) (*domain.User, *errors.RestErr) {
	var resp domain.User

	res := r.DB.Model(&domain.User{}).Where("id = ?", userID).First(&resp)

	if res.RowsAffected == 0 {
		logger.Error("error occurred when getting user by user id", res.Error)
		return nil, errors.NewNotFoundError(errors.ErrRecordNotFound)
	}

	if res.Error != nil {
		logger.Error("error occurred when getting user by user id", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return &resp, nil
}

func (r *users) Update(user *domain.User) *errors.RestErr {
	res := r.DB.Model(&domain.User{}).Omit("password").Where("id = ?", user.ID).Updates(&user)

	if res.Error != nil {
		logger.Error("error occurred when updating user by user id", res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *users) UpdatePassword(userID uint, updateValues map[string]interface{}) *errors.RestErr {
	res := r.DB.Model(&domain.User{}).Where("id = ?", userID).Updates(&updateValues)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("updating user by user id"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *users) GetUserByAppKey(appKey string) (*domain.User, *errors.RestErr) {
	var resp domain.User

	res := r.DB.Model(&domain.User{}).Where("app_key = ?", appKey).First(&resp)

	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError("no user found")
	}

	if res.Error != nil {
		logger.Error("error occurred when getting user by app key", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return &resp, nil
}

func (r *users) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	res := r.DB.Model(&domain.User{}).Where("email = ?", email).Find(&user)
	if res.RowsAffected == 0 {
		logger.Error("no user found by this email", res.Error)
		return nil, errors.NewError(errors.ErrRecordNotFound)
	}
	if res.Error != nil {
		logger.Error("error occurred when trying to get user by email", res.Error)
		return nil, errors.NewError(errors.ErrSomethingWentWrong)
	}

	return user, nil
}

func (r *users) SetLastLoginAt(user *domain.User) error {
	lastLoginAt := time.Now().UTC()

	err := r.DB.Model(&user).Update("last_login_at", lastLoginAt).Error

	if err != nil {
		logger.Error(err.Error(), err)
		return err
	}

	return nil
}

func (r *users) HasRole(userID, roleID uint) bool {
	var count int64
	count = 0

	r.DB.Model(&domain.User{}).
		Where("id = ? AND role_id = ?", userID, roleID).
		Count(&count)

	return count > 0
}

func (r *users) ResetPassword(userID int, hashedPass []byte) error {
	err := r.DB.Model(&domain.User{}).
		Where("id = ?", userID).
		Update("password", hashedPass).
		Error

	if err != nil {
		logger.Error("error occur when reset password", err)
		return err
	}

	return nil
}

func (r *users) GetTokenUser(id uint) (*domain.VerifyTokenResp, *errors.RestErr) {
	tempUser := &domain.TempVerifyTokenResp{}
	var vtUser domain.VerifyTokenResp

	query := r.tokenUserFetchQuery()

	res := query.Where("users.id = ?", id).Find(&tempUser)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("get token user"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	err := methodsutil.StructToStruct(tempUser, &vtUser.BaseVerifyTokenResp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	vtUser.Permissions = strings.Split(tempUser.Permissions, ",")

	return &vtUser, nil
}

func (r *users) GetUserWithPermissions(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr) {
	var intUser domain.IntermediateUserWithPermissions
	var userWithParams domain.UserWithPerms

	sections := `
		users.*,
		roles.name role_name
	`
	if withPermission {
		sections += ",GROUP_CONCAT(DISTINCT permissions.name) AS permissions"
	}

	query := r.DB.Model(&domain.User{}).
		Select(sections).
		Joins("LEFT JOIN roles ON users.role_id = roles.id").
		Where("users.deleted_at IS NULL")

	if withPermission {
		query = query.
			Joins("JOIN role_permissions ON users.role_id = role_permissions.role_id").
			Joins("JOIN permissions ON role_permissions.permission_id = permissions.id")
	}

	query.Group("users.id")

	res := query.Where("users.id = ?", userID).Find(&intUser)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("getting user with permission"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	err := methodsutil.StructToStruct(intUser, &userWithParams.User)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	userWithParams.RoleName = intUser.RoleName

	if withPermission {
		userWithParams.Permissions = strings.Split(intUser.Permissions, ",")
	}

	return &userWithParams, nil
}

func (r *users) GetUserRankListByPoint() ([]*domain.User, *errors.RestErr) {
	users := []*domain.User{}

	err := r.DB.Model(&domain.User{}).Order("total_point desc").Find(&users).Error
	if err != nil {
		logger.Error("error occurred when getting all users", err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return users, nil
}

func (r *users) tokenUserFetchQuery() *gorm.DB {
	selections := `
		users.id,
		users.first_name,
		users.last_name,
		users.email,
		users.phone,
		users.profile_pic,
		companies.business_id,
		businesses.name business_name,
		companies.id company_id,
		companies.name company_name,
		(
			CASE
				WHEN 1 IN (GROUP_CONCAT(DISTINCT users.role_id)) THEN 1 ELSE 0
			END
		) AS admin,
		(
			CASE
				WHEN 3 IN (GROUP_CONCAT(DISTINCT users.role_id)) THEN 1 ELSE 0
			END
		) AS super_admin,
		GROUP_CONCAT(DISTINCT permissions.name) AS permissions
	`

	return r.DB.Table("users").
		Select(selections).
		Joins("LEFT JOIN companies ON users.company_id = companies.id").
		Joins("LEFT JOIN businesses ON companies.business_id = businesses.id").
		Joins("JOIN roles ON users.role_id = roles.id").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("users.deleted_at IS NULL").
		Group("users.id")
}

// SaveCommitments commitments
func (r *users) SaveCommitments(commitments domain.Commitments) (*domain.Commitments, *errors.RestErr) {

	usr, _ := r.GetUserByID(commitments.DoctorId)
	upErr := r.DB.Model(&domain.User{}).Where("id = ?", usr.ID).Updates(map[string]interface{}{
		"total_point": usr.TotalPoint + commitments.Point,
	}).Error

	if upErr != nil {
		logger.Error("error occurred when updating user by user id", upErr)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	gUsr, _ := r.GetUserByID(commitments.DoctorId)
	commitments.Point = gUsr.TotalPoint
	res := r.DB.Model(&domain.Commitments{}).Create(&commitments)

	if res.Error != nil {
		logger.Error("error occurred when create commitments", res.Error)
		restErr := fmt.Sprintf("%s", res.Error)
		return nil, errors.NewInternalServerError(fmt.Sprintf("%s", errors.NewError(restErr)))
	}

	return &commitments, nil
}

func (r *users) AllCommitments(cid uint) ([]*domain.Commitments, *errors.RestErr) {
	commitments := []*domain.Commitments{}

	err := r.DB.Model(&domain.Commitments{}).Where("doctor_id = ?", cid).Find(&commitments).Error
	if err != nil {
		logger.Error("error occurred when getting all Commitments", err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return commitments, nil
}

func (r *users) DeleteCommitments(cid uint) *errors.RestErr {
	err := r.DB.Model(&domain.Commitments{}).Where("id = ?", cid).Delete(&domain.Commitments{}).Error
	if err != nil {
		logger.Error("error occurred when deleting  Commitments", err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return nil
}

// Help  stuff
func (r *users) SaveHelp(help *domain.Help) (*domain.Help, *errors.RestErr) {
	res := r.DB.Model(&domain.Help{}).Create(&help)

	if res.Error != nil {
		logger.Error("error occurred when create user", res.Error)
		restErr := fmt.Sprintf("%s", res.Error)
		return nil, errors.NewInternalServerError(fmt.Sprintf("%s", errors.NewError(restErr)))
	}

	return help, nil
}

func (r *users) UpdateHelp(help *domain.Help) *errors.RestErr {
	res := r.DB.Model(&domain.User{}).Where("id = ?", help.ID).Updates(&help)

	if res.Error != nil {
		logger.Error("error occurred when updating user by user id", res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

// Place stuff
func (r *users) SavePlace(place domain.Place) (*domain.Place, *errors.RestErr) {
	res := r.DB.Model(&domain.Place{}).Create(&place)

	if res.Error != nil {
		logger.Error("error occurred when create user", res.Error)
		restErr := fmt.Sprintf("%s", res.Error)
		return nil, errors.NewInternalServerError(fmt.Sprintf("%s", errors.NewError(restErr)))
	}

	return &place, nil
}

func (r *users) AllPlaces() ([]*domain.Place, *errors.RestErr) {
	places := []*domain.Place{}

	err := r.DB.Model(&domain.Place{}).Find(&places).Error
	if err != nil {
		logger.Error("error occurred when getting all Commitments", err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return places, nil
}

func (r *users) DeletePLace(pid uint) *errors.RestErr {
	err := r.DB.Model(&domain.Place{}).Where("id = ?", pid).Delete(&domain.Place{}).Error
	if err != nil {
		logger.Error("error occurred when deleting  Commitments", err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return nil
}
