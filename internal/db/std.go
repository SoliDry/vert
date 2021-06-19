package db

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/SoliDry/vert"
	"github.com/arthurkushman/pgo"
)

// std dirs
const (
	vertRootLibDirLinux = "/var/lib/vert"
	vertRootLibDirMacOS = "/usr/local/vert/data"
	infoDbPath          = "/information_schema"
)

// std tables it is a path representing table in which will be written discrete files with fields
const (
	sysTblUser       = "user"
	sysTblDb         = "db"
	sysTblIndexStats = "index_stats"
	sysTblStats      = "table_stats"
)

// supported OSs
const (
	osDarwin = "darwin"
	osLinux  = "linux"
)

const (
	tblStructFileFormat = ".fmt"
	tblDataFileFormat   = ".db"
)

var stdTables = map[string]map[string][]string{
	sysTblUser: {
		"Host":        []string{"CHAR(60)", "NOT NULL", "PRIMARY KEY"}, // Host (together with User makes up the unique identifier for this account.
		"User":        []string{"CHAR(80)", "NOT NULL", "PRIMARY KEY"},
		"Password":    []string{"CHAR(41)", "NOT NULL"},
		"Select_priv": []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Insert_priv": []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Update_priv": []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Delete_priv": []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Create_priv": []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Drop_priv":   []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Index_priv":  []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
	},
	sysTblDb: {
		"Host":        []string{"CHAR(60)", "NOT NULL", "PRIMARY KEY"}, // Host (together with User makes up the unique identifier for this account.
		"Db":          []string{"CHAR(64)", "NOT NULL", "PRIMARY KEY"},
		"User":        []string{"CHAR(80)", "NOT NULL", "PRIMARY KEY"},
		"Select_priv": []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Insert_priv": []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Update_priv": []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Delete_priv": []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Create_priv": []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Drop_priv":   []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Index_priv":  []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Grant_priv":  []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
		"Alter_priv":  []string{"ENUM('N', 'Y')", "NOT NULL", "DEFAULT 'N'"},
	},
	sysTblIndexStats: {
		"db_name":          []string{"VARCHAR(64)", "NOT NULL", "PRIMARY KEY"}, // Host (together with User makes up the unique identifier for this account.
		"table_name":       []string{"VARCHAR(64)", "NOT NULL", "PRIMARY KEY"},
		"index_name":       []string{"VARCHAR(64)", "NOT NULL", "PRIMARY KEY"},
		"avg_frequency":    []string{"DECIMAL(12, 4)", "NULL", "DEFAULT NULL"},
		"last_update":      []string{"TIMESTAMP", "NOT NULL", "DEFAULT CURRENT_TIMESTAMP"},
		"stat_value":       []string{"BIGINT UNSIGNED", "NOT NULL", "DEFAULT 0"},
		"sample_size":      []string{"BIGINT UNSIGNED", "NULL", "DEFAULT NULL"},
		"stat_description": []string{"VARCHAR(1024)", "NOT NULL", "DEFAULT ''"},
	},
	sysTblStats: {
		"db_name":                  []string{"VARCHAR(64)", "NOT NULL", "PRIMARY KEY"}, // Host (together with User makes up the unique identifier for this account.
		"table_name":               []string{"VARCHAR(64)", "NOT NULL", "PRIMARY KEY"},
		"last_update":              []string{"TIMESTAMP", "NOT NULL", "DEFAULT CURRENT_TIMESTAMP"},
		"n_rows":                   []string{"BIGINT UNSIGNED", "NOT NULL", "DEFAULT 0"},
		"clustered_index_size":     []string{"BIGINT UNSIGNED", "NOT NULL", "DEFAULT 0"},
		"sum_of_other_index_sizes": []string{"BIGINT UNSIGNED", "NOT NULL", "DEFAULT 0"},
	},
}

// CreateSysFilesAndTables creates system tables (if not exist) like user, db, slow_log etc
func CreateSysFilesAndTables() {
	if pgo.FileExists(vertRootLibDirLinux+infoDbPath) == false {
		if runtime.GOOS == osLinux {
			createSysFiles(vertRootLibDirLinux)
		}

		if runtime.GOOS == osDarwin {
			createSysFiles(vertRootLibDirMacOS)
		}
	}
}

func createSysFiles(rootPath string) {
	if pgo.FileExists(rootPath) == false {
		vert.MakeDir(rootPath, 755)
	}

	vert.MakeDir(rootPath+infoDbPath, 755)

	// creating user table
	vert.MakeDir(rootPath+infoDbPath+sysTblUser, 755)
	createColumns(sysTblUser)

	// creating db table
	vert.MakeDir(rootPath+infoDbPath+sysTblDb, 755)
	createColumns(sysTblDb)
}

func createColumns(sysTbl string) {
	userMap, ok := stdTables[sysTbl]
	if ok {
		for colN, opts := range userMap {
			f, err := os.OpenFile(colN+tblStructFileFormat, os.O_RDWR|os.O_CREATE, 755)
			if err != nil {
				log.Fatal(fmt.Errorf("could not create file: %w", err))
			}

			tblContent := ""
			for i, opt := range opts {
				if i == 0 {
					tblContent = opt
				} else {
					tblContent += "|" + opt
				}
			}

			f.Write([]byte(tblContent + "\n"))
			f.Close()
		}
	}
}
