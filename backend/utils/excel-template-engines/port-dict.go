package excel_template_engines

import (
	"strings"
	"sync"

	"super-supply-chain/models"
)

const portDictType = "港口字典"

var (
	portDictMu    sync.RWMutex
	portDictCache map[string]PortInfo
)

// InvalidatePortDictCache clears the in-memory 港口字典 cache (tests / after admin edits).
func InvalidatePortDictCache() {
	portDictMu.Lock()
	portDictCache = nil
	portDictMu.Unlock()
}

// LoadPortDict reads type=港口字典 into PortInfo keyed by arrival_port text.
func LoadPortDict() map[string]PortInfo {
	portDictMu.RLock()
	if portDictCache != nil {
		out := portDictCache
		portDictMu.RUnlock()
		return out
	}
	portDictMu.RUnlock()

	portDictMu.Lock()
	defer portDictMu.Unlock()
	if portDictCache != nil {
		return portDictCache
	}

	var rows []models.BaseDict
	q := models.DB.Model(&models.BaseDict{}).Where("type = ?", portDictType).Find(&rows)
	m := make(map[string]PortInfo, len(rows))
	if q.Error == nil {
		for _, d := range rows {
			key := strings.TrimSpace(d.Key)
			if key == "" {
				continue
			}
			addr := strings.TrimSpace(d.Value)
			portName := strings.TrimSpace(d.PortName)
			if portName == "" && addr != "" {
				portName = addr + "口岸"
			}
			m[key] = PortInfo{
				PortName: portName,
				Addr:     addr,
				ExtraPay: strings.TrimSpace(d.ExtraPay),
			}
		}
	}
	portDictCache = m
	return portDictCache
}

// GetPortInfo returns 港口字典 entry for an arrival_port string.
func GetPortInfo(port string) PortInfo {
	port = strings.TrimSpace(port)
	if port == "" {
		return PortInfo{}
	}
	return LoadPortDict()[port]
}
