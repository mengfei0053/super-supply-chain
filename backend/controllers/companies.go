package controllers

import (
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

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
