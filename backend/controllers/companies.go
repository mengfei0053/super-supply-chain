package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"super-supply-chain/models"
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

const (
	companyNameMax       = 255
	companyAliasMax      = 255
	companyTargetAddrMax = 255
	companyCreditCodeMax = 100
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
}

type companyUpdateRequest struct {
	Name       *string `json:"name"`
	Alias      *string `json:"alias"`
	TargetAddr *string `json:"target_addr"`
}

// CreateCompany inserts one base_companies_infos row.
// name and unified_social_credit_code are required by the existing model.
// alias and target_addr are optional. Other columns are left empty.
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
	alias := ""
	if req.Alias != nil {
		alias, errMsg = normalizeOptionalCompanyField(*req.Alias, "alias", companyAliasMax)
		if errMsg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
	}
	targetAddr := ""
	if req.TargetAddr != nil {
		targetAddr, errMsg = normalizeOptionalCompanyField(*req.TargetAddr, "target_addr", companyTargetAddrMax)
		if errMsg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
	}
	code, errMsg := normalizeRequiredCompanyField(req.UnifiedSocialCreditCode, "unified_social_credit_code", companyCreditCodeMax)
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
			) VALUES (?, ?, ?, ?, '', '', '', '', '', '', ?, ?, NULL)`,
			name, alias, targetAddr, code, now, now,
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

// UpdateCompany changes name, alias, and/or target_addr on one active company.
// Omitted JSON fields stay as stored. An empty alias or target_addr clears that column.
// Credit code, address, bank, and phone are not modified.
func UpdateCompany(c *gin.Context) {
	id64, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	id := uint(id64)

	var req companyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name == nil && req.Alias == nil && req.TargetAddr == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, alias, or target_addr is required"})
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
	var alias string
	if req.Alias != nil {
		var errMsg string
		alias, errMsg = normalizeOptionalCompanyField(*req.Alias, "alias", companyAliasMax)
		if errMsg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
	}
	var targetAddr string
	if req.TargetAddr != nil {
		var errMsg string
		targetAddr, errMsg = normalizeOptionalCompanyField(*req.TargetAddr, "target_addr", companyTargetAddrMax)
		if errMsg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
	}

	var updated companyInfoRow
	err = models.DB.Transaction(func(tx *gorm.DB) error {
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

		sets := make([]string, 0, 4)
		args := make([]any, 0, 5)
		if req.Name != nil {
			sets = append(sets, "name = ?")
			args = append(args, name)
		}
		if req.Alias != nil {
			sets = append(sets, "alias = ?")
			args = append(args, alias)
		}
		if req.TargetAddr != nil {
			sets = append(sets, "target_addr = ?")
			args = append(args, targetAddr)
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
