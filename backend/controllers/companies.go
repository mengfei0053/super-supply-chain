package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"super-supply-chain/models"
	"super-supply-chain/utils"
)

const companyInfoColumns = "id, created_at, updated_at, deleted_at, name, addr_country, addr_province, addr_city, addr_street, unified_social_credit_code, bank_code, phone_num, alias, target_addr"

type companyInfoRow struct {
	ID                      uint       `gorm:"column:id" json:"id"`
	CreatedAt               time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt               time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt               *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Name                    string     `gorm:"column:name" json:"name"`
	AddrCountry             string     `gorm:"column:addr_country" json:"addr_country"`
	AddrProvince            string     `gorm:"column:addr_province" json:"addr_province"`
	AddrCity                string     `gorm:"column:addr_city" json:"addr_city"`
	AddrStreet              string     `gorm:"column:addr_street" json:"addr_street"`
	UnifiedSocialCreditCode string     `gorm:"column:unified_social_credit_code" json:"unified_social_credit_code"`
	BankCode                string     `gorm:"column:bank_code" json:"bank_code"`
	PhoneNum                string     `gorm:"column:phone_num" json:"phone_num"`
	Alias                   string     `gorm:"column:alias" json:"alias"`
	TargetAddr              string     `gorm:"column:target_addr" json:"target_addr"`
}

const (
	companyNameMax       = 255
	companyAliasMax      = 255
	companyTargetAddrMax = 255
	companyCreditCodeMax = 100
	companyAddrMax       = 255
	companyAddrPartMax   = 100
	companyBankCodeMax   = 100
	companyPhoneMax      = 20
)

var (
	errCompanyNameExists = errors.New("name already exists")
	errCompanyCodeExists = errors.New("unified_social_credit_code already exists")
	errCompanyNotFound   = errors.New("company not found")
)

type companyCreateRequest struct {
	Name                    string  `json:"name"`
	Alias                   *string `json:"alias"`
	TargetAddr              *string `json:"target_addr"`
	UnifiedSocialCreditCode string  `json:"unified_social_credit_code"`
	AddrCountry             *string `json:"addr_country"`
	AddrProvince            *string `json:"addr_province"`
	AddrCity                *string `json:"addr_city"`
	AddrStreet              *string `json:"addr_street"`
	BankCode                *string `json:"bank_code"`
	PhoneNum                *string `json:"phone_num"`
}

type companyUpdateRequest struct {
	Name         *string `json:"name"`
	Alias        *string `json:"alias"`
	TargetAddr   *string `json:"target_addr"`
	AddrCountry  *string `json:"addr_country"`
	AddrProvince *string `json:"addr_province"`
	AddrCity     *string `json:"addr_city"`
	AddrStreet   *string `json:"addr_street"`
	BankCode     *string `json:"bank_code"`
	PhoneNum     *string `json:"phone_num"`
}

// CompaniesIndex handles GET /companies.
// With keyword → keyword search (MCP / curl).
// Without keyword → react-admin list (filter/sort/range + Content-Range).
func CompaniesIndex(c *gin.Context) {
	if strings.TrimSpace(c.Query("keyword")) != "" {
		SearchCompanies(c)
		return
	}
	ListCompanies(c)
}

// SearchCompanies lists base_companies_infos rows whose name or alias contains keyword.
// Rows with deleted_at set are omitted unless includeDeleted is true.
func SearchCompanies(c *gin.Context) {
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}
	if utf8.RuneCountInString(keyword) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is too long"})
		return
	}

	pattern := likeContains(keyword)
	escape := likeEscapeSQL()
	rows := make([]companyInfoRow, 0)
	query := models.DB.Table("base_companies_infos").
		Select(companyInfoColumns).
		Where("(name LIKE ? "+escape+" OR alias LIKE ? "+escape+")", pattern, pattern)
	if !includeDeletedQuery(c) {
		query = query.Where("deleted_at IS NULL")
	}
	if err := query.Order("id asc").Scan(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// ListCompanies returns a paginated react-admin list of active companies.
// Optional filter.q / filter.keyword / filter.name matches name or alias.
func ListCompanies(c *gin.Context) {
	listQuery, err := utils.GetListQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	keyword := companyListFilterKeyword(c)
	if utf8.RuneCountInString(keyword) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is too long"})
		return
	}

	applyFilters := func(db *gorm.DB) *gorm.DB {
		db = db.Table("base_companies_infos").Where("deleted_at IS NULL")
		if keyword != "" {
			pattern := likeContains(keyword)
			escape := likeEscapeSQL()
			db = db.Where("(name LIKE ? "+escape+" OR alias LIKE ? "+escape+")", pattern, pattern)
		}
		return db
	}

	var total int64
	if err := applyFilters(models.DB).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows := make([]companyInfoRow, 0)
	if err := applyFilters(models.DB).Select(companyInfoColumns).
		Order("id asc").
		Limit(listQuery.Limit).
		Offset(listQuery.Offset).
		Scan(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SetContentRange(c, total)
	c.JSON(http.StatusOK, rows)
}

// GetCompany returns one active company by id.
func GetCompany(c *gin.Context) {
	id, ok := parseCompanyIDParam(c)
	if !ok {
		return
	}
	row, err := loadActiveCompanyByID(models.DB, id)
	if err != nil {
		writeCompanyResult(c, err, companyInfoRow{})
		return
	}
	c.JSON(http.StatusOK, row)
}

// CreateCompany inserts one base_companies_infos row.
// name and unified_social_credit_code are required by the existing model.
// alias, target_addr, address, bank_code, and phone_num are optional.
func CreateCompany(c *gin.Context) {
	var req companyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name, errMsg := normalizeRequiredCompanyField(req.Name, "name", companyNameMax)
	if errMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	code, errMsg := normalizeRequiredCompanyField(req.UnifiedSocialCreditCode, "unified_social_credit_code", companyCreditCodeMax)
	if errMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	fields, errMsg := normalizeOptionalCompanyBasics(
		req.Alias, req.TargetAddr,
		req.AddrCountry, req.AddrProvince, req.AddrCity, req.AddrStreet,
		req.BankCode, req.PhoneNum,
	)
	if errMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	var created companyInfoRow
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		if id, err := companyIDBy(tx, "name", name, 0); err != nil {
			return err
		} else if id != 0 {
			return errCompanyNameExists
		}
		if id, err := companyIDBy(tx, "unified_social_credit_code", code, 0); err != nil {
			return err
		} else if id != 0 {
			return errCompanyCodeExists
		}

		now := time.Now()
		if err := tx.Exec(
			`INSERT INTO base_companies_infos (
				name, alias, target_addr, unified_social_credit_code,
				addr_country, addr_province, addr_city, addr_street,
				bank_code, phone_num, created_at, updated_at, deleted_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL)`,
			name, fields.alias, fields.targetAddr, code,
			fields.addrCountry, fields.addrProvince, fields.addrCity, fields.addrStreet,
			fields.bankCode, fields.phoneNum, now, now,
		).Error; err != nil {
			return err
		}

		row, err := loadActiveCompanyByName(tx, name)
		if err != nil {
			return err
		}
		created = row
		return nil
	})
	writeCompanyResult(c, err, created)
}

// UpdateCompany changes basic-info fields on one active company.
// Omitted JSON fields stay as stored. An empty optional string clears that column.
// unified_social_credit_code is not modified.
func UpdateCompany(c *gin.Context) {
	id, ok := parseCompanyIDParam(c)
	if !ok {
		return
	}

	var req companyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == nil && req.Alias == nil && req.TargetAddr == nil &&
		req.AddrCountry == nil && req.AddrProvince == nil && req.AddrCity == nil &&
		req.AddrStreet == nil && req.BankCode == nil && req.PhoneNum == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one updatable field is required"})
		return
	}

	var name string
	if req.Name != nil {
		var errMsg string
		name, errMsg = normalizeRequiredCompanyField(*req.Name, "name", companyNameMax)
		if errMsg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
	}
	fields, errMsg := normalizeOptionalCompanyBasics(
		req.Alias, req.TargetAddr,
		req.AddrCountry, req.AddrProvince, req.AddrCity, req.AddrStreet,
		req.BankCode, req.PhoneNum,
	)
	if errMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	var updated companyInfoRow
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		if _, err := loadActiveCompanyByID(tx, id); err != nil {
			return err
		}
		if req.Name != nil {
			otherID, err := companyIDBy(tx, "name", name, id)
			if err != nil {
				return err
			}
			if otherID != 0 {
				return errCompanyNameExists
			}
		}

		sets := make([]string, 0, 10)
		args := make([]any, 0, 11)
		if req.Name != nil {
			sets = append(sets, "name = ?")
			args = append(args, name)
		}
		if req.Alias != nil {
			sets = append(sets, "alias = ?")
			args = append(args, fields.alias)
		}
		if req.TargetAddr != nil {
			sets = append(sets, "target_addr = ?")
			args = append(args, fields.targetAddr)
		}
		if req.AddrCountry != nil {
			sets = append(sets, "addr_country = ?")
			args = append(args, fields.addrCountry)
		}
		if req.AddrProvince != nil {
			sets = append(sets, "addr_province = ?")
			args = append(args, fields.addrProvince)
		}
		if req.AddrCity != nil {
			sets = append(sets, "addr_city = ?")
			args = append(args, fields.addrCity)
		}
		if req.AddrStreet != nil {
			sets = append(sets, "addr_street = ?")
			args = append(args, fields.addrStreet)
		}
		if req.BankCode != nil {
			sets = append(sets, "bank_code = ?")
			args = append(args, fields.bankCode)
		}
		if req.PhoneNum != nil {
			sets = append(sets, "phone_num = ?")
			args = append(args, fields.phoneNum)
		}
		sets = append(sets, "updated_at = ?")
		args = append(args, time.Now())
		args = append(args, id)

		if err := tx.Exec(
			"UPDATE base_companies_infos SET "+strings.Join(sets, ", ")+" WHERE id = ? AND deleted_at IS NULL",
			args...,
		).Error; err != nil {
			return err
		}
		row, err := loadActiveCompanyByID(tx, id)
		if err != nil {
			return err
		}
		updated = row
		return nil
	})
	writeCompanyResult(c, err, updated)
}

// DeleteCompany soft-deletes one active company by setting deleted_at.
func DeleteCompany(c *gin.Context) {
	id, ok := parseCompanyIDParam(c)
	if !ok {
		return
	}

	now := time.Now()
	result := models.DB.Exec(
		`UPDATE base_companies_infos SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`,
		now, now, id,
	)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": errCompanyNotFound.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Delete company successfully", "id": id})
}

func writeCompanyResult(c *gin.Context, err error, row companyInfoRow) {
	switch {
	case err == nil:
		c.JSON(http.StatusOK, row)
	case errors.Is(err, errCompanyNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": errCompanyNotFound.Error()})
	case errors.Is(err, errCompanyNameExists):
		c.JSON(http.StatusConflict, gin.H{"error": errCompanyNameExists.Error()})
	case errors.Is(err, errCompanyCodeExists):
		c.JSON(http.StatusConflict, gin.H{"error": errCompanyCodeExists.Error()})
	case isDuplicateKey(err):
		c.JSON(http.StatusConflict, gin.H{"error": "name or unified_social_credit_code already exists"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

type companyOptionalBasics struct {
	alias        string
	targetAddr   string
	addrCountry  string
	addrProvince string
	addrCity     string
	addrStreet   string
	bankCode     string
	phoneNum     string
}

func normalizeOptionalCompanyBasics(
	alias, targetAddr, addrCountry, addrProvince, addrCity, addrStreet, bankCode, phoneNum *string,
) (companyOptionalBasics, string) {
	out := companyOptionalBasics{}
	var errMsg string
	if alias != nil {
		out.alias, errMsg = normalizeOptionalCompanyField(*alias, "alias", companyAliasMax)
		if errMsg != "" {
			return out, errMsg
		}
	}
	if targetAddr != nil {
		out.targetAddr, errMsg = normalizeOptionalCompanyField(*targetAddr, "target_addr", companyTargetAddrMax)
		if errMsg != "" {
			return out, errMsg
		}
	}
	if addrCountry != nil {
		out.addrCountry, errMsg = normalizeOptionalCompanyField(*addrCountry, "addr_country", companyAddrPartMax)
		if errMsg != "" {
			return out, errMsg
		}
	}
	if addrProvince != nil {
		out.addrProvince, errMsg = normalizeOptionalCompanyField(*addrProvince, "addr_province", companyAddrPartMax)
		if errMsg != "" {
			return out, errMsg
		}
	}
	if addrCity != nil {
		out.addrCity, errMsg = normalizeOptionalCompanyField(*addrCity, "addr_city", companyAddrPartMax)
		if errMsg != "" {
			return out, errMsg
		}
	}
	if addrStreet != nil {
		out.addrStreet, errMsg = normalizeOptionalCompanyField(*addrStreet, "addr_street", companyAddrMax)
		if errMsg != "" {
			return out, errMsg
		}
	}
	if bankCode != nil {
		out.bankCode, errMsg = normalizeOptionalCompanyField(*bankCode, "bank_code", companyBankCodeMax)
		if errMsg != "" {
			return out, errMsg
		}
	}
	if phoneNum != nil {
		out.phoneNum, errMsg = normalizeOptionalCompanyField(*phoneNum, "phone_num", companyPhoneMax)
		if errMsg != "" {
			return out, errMsg
		}
	}
	return out, ""
}

func normalizeRequiredCompanyField(raw string, field string, max int) (string, string) {
	value, errMsg := normalizeOptionalCompanyField(raw, field, max)
	if errMsg != "" {
		return "", errMsg
	}
	if value == "" {
		return "", field + " is required"
	}
	return value, ""
}

func normalizeOptionalCompanyField(raw string, field string, max int) (string, string) {
	value := strings.TrimSpace(raw)
	if utf8.RuneCountInString(value) > max {
		return "", field + " is too long"
	}
	return value, ""
}

func parseCompanyIDParam(c *gin.Context) (uint, bool) {
	id64, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(id64), true
}

func companyListFilterKeyword(c *gin.Context) string {
	raw := strings.TrimSpace(c.Query("filter"))
	if raw == "" {
		return ""
	}
	var filter map[string]any
	if err := json.Unmarshal([]byte(raw), &filter); err != nil {
		return ""
	}
	for _, key := range []string{"q", "keyword", "name"} {
		if value, ok := filter[key]; ok {
			switch v := value.(type) {
			case string:
				return strings.TrimSpace(v)
			}
		}
	}
	return ""
}

func companyIDBy(db *gorm.DB, column, value string, exceptID uint) (uint, error) {
	switch column {
	case "name", "unified_social_credit_code":
	default:
		return 0, errors.New("unsupported company column")
	}
	query := db.Table("base_companies_infos").Select("id").Where(column+" = ?", value)
	if exceptID != 0 {
		query = query.Where("id <> ?", exceptID)
	}
	var rows []struct {
		ID uint `gorm:"column:id"`
	}
	if err := query.Limit(1).Find(&rows).Error; err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}
	return rows[0].ID, nil
}

func loadActiveCompanyByID(db *gorm.DB, id uint) (companyInfoRow, error) {
	var rows []companyInfoRow
	err := db.Table("base_companies_infos").
		Select(companyInfoColumns).
		Where("id = ? AND deleted_at IS NULL", id).
		Limit(1).
		Scan(&rows).Error
	if err != nil {
		return companyInfoRow{}, err
	}
	if len(rows) == 0 {
		return companyInfoRow{}, errCompanyNotFound
	}
	return rows[0], nil
}

func loadActiveCompanyByName(db *gorm.DB, name string) (companyInfoRow, error) {
	var rows []companyInfoRow
	err := db.Table("base_companies_infos").
		Select(companyInfoColumns).
		Where("name = ? AND deleted_at IS NULL", name).
		Limit(1).
		Scan(&rows).Error
	if err != nil {
		return companyInfoRow{}, err
	}
	if len(rows) == 0 {
		return companyInfoRow{}, errCompanyNotFound
	}
	return rows[0], nil
}

func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") ||
		strings.Contains(msg, "duplicate entry") ||
		strings.Contains(msg, "error 1062")
}

func includeDeletedQuery(c *gin.Context) bool {
	switch strings.ToLower(strings.TrimSpace(c.Query("includeDeleted"))) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}

func likeContains(keyword string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + replacer.Replace(keyword) + "%"
}

// MySQL string literals treat backslash as an escape. SQLite does not.
func likeEscapeSQL() string {
	if models.DB != nil && models.DB.Dialector.Name() == "mysql" {
		return `ESCAPE '\\'`
	}
	return `ESCAPE '\'`
}
