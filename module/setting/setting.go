package setting

import (
	"net/mail"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/MessageDream/goby/module/mailer"
	"github.com/MessageDream/goby/module/storage"

	"github.com/Unknwon/com"
	"github.com/go-macaron/session"
	"gopkg.in/ini.v1"

	"github.com/MessageDream/goby/model"
	log "gopkg.in/clog.v1"
)

type Scheme string

const (
	SCHEME_HTTP  Scheme = "http"
	SCHEME_HTTPS Scheme = "https"
)

var (
	// App settings.
	AppVer      string
	AppName     string
	AppLogo     string
	AppURL      string
	AppSubURL   string
	AppDataPath string

	// Server settings.
	Protocol           Scheme
	Domain             string
	HTTPAddr, HTTPPort string
	LocalURL           string
	SSHPort            int
	OfflineMode        bool
	DisableRouterLog   bool
	CertFile, KeyFile  string
	StaticRootPath     string
	EnableGzip         bool

	// Attachment settings
	AttachmentPath         string
	AttachmentAllowedTypes string
	AttachmentMaxSize      int64
	AttachmentMaxFiles     int
	AttachmentEnabled      bool

	// Database settings
	UseSQLite3    bool
	UseMySQL      bool
	UsePostgreSQL bool
	UseMSSQL      bool
	UseTiDB       bool

	// Security settings.
	InstallLock          bool
	SecretKey            string
	LogInRememberDays    int
	CookieUserName       string
	CookieRememberName   string
	CookieSecure         bool
	ReverseProxyAuthUser string

	// Log settings.
	LogRootPath string
	LogModes    []string
	LogConfigs  []interface{}

	// Time settings.
	TimeFormat string

	// Cache settings.
	CacheAdapter                  string
	CacheInternal                 int
	CacheConn                     string
	CacheClientCheckUpdateTimeOut int64
	CacheTokenTimeOut             int64

	EnableRedis    bool
	EnableMemcache bool

	// Session settings.
	SessionConfig  session.Options
	CSRFCookieName = "_csrf"

	// Global setting objects.
	Cfg          *ini.File
	ConfRootPath string
	CustomPath   string // Custom directory path.
	CustomConf   string
	ProdMode     bool
	RunUser      string

	HasRobotsTxt bool
)

func ExecPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	p, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return p, nil
}

func WorkDir() (string, error) {
	execPath, err := ExecPath()
	return path.Dir(strings.Replace(execPath, "\\", "/", -1)), err
}

func InitConfig() {
	workDir, err := WorkDir()
	if err != nil {
		log.Fatal(2, "Fail to get work directory: %v", err)
	}

	Cfg, err = ini.Load(path.Join("conf/app.ini"))
	if err != nil {
		log.Fatal(2, "Fail to parse 'conf/app.ini': %v", err)
	}

	CustomPath = os.Getenv("GOBY_CUSTOM")
	if len(CustomPath) == 0 {
		CustomPath = workDir + "/custom"
	}

	if len(CustomConf) == 0 {
		CustomConf = CustomPath + "/conf/app.ini"
	}

	if com.IsFile(CustomConf) {
		if err = Cfg.Append(CustomConf); err != nil {
			log.Fatal(2, "Fail to load custom conf '%s': %v", CustomConf, err)
		}
	} else {
		log.Warn("Custom config '%s' not found, ignore this if you're running first time", CustomConf)
	}
	Cfg.NameMapper = ini.AllCapsUnderscore

	homeDir, err := com.HomeDir()
	if err != nil {
		log.Fatal(2, "Fail to get home directory: %v", err)
	}
	homeDir = strings.Replace(homeDir, "\\", "/", -1)

	LogRootPath = Cfg.Section("log").Key("ROOT_PATH").MustString(path.Join(workDir, "log"))

	sec := Cfg.Section("server")
	AppName = Cfg.Section("").Key("APP_NAME").MustString("Goby")
	AppURL = sec.Key("ROOT_URL").MustString("http://localhost:3000/")
	if AppURL[len(AppURL)-1] != '/' {
		AppURL += "/"
	}

	// Check if has app SubURL.
	url, err := url.Parse(AppURL)
	if err != nil {
		log.Fatal(2, "Invalid ROOT_URL '%s': %s", AppURL, err)
	}
	// SubURL should start with '/' and end without '/', such as '/{subpath}'.
	// This value is empty if site does not have sub-url.
	AppSubURL = strings.TrimSuffix(url.Path, "/")

	Protocol = SCHEME_HTTP
	if sec.Key("PROTOCOL").String() == "https" {
		Protocol = SCHEME_HTTPS
		CertFile = sec.Key("CERT_FILE").String()
		KeyFile = sec.Key("KEY_FILE").String()
	}
	Domain = sec.Key("DOMAIN").MustString("localhost")
	HTTPAddr = sec.Key("HTTP_ADDR").MustString("0.0.0.0")
	HTTPPort = sec.Key("HTTP_PORT").MustString("3000")
	LocalURL = sec.Key("LOCAL_ROOT_URL").MustString(string(Protocol) + "://localhost:" + HTTPPort + "/")
	OfflineMode = sec.Key("OFFLINE_MODE").MustBool()
	DisableRouterLog = sec.Key("DISABLE_ROUTER_LOG").MustBool()
	StaticRootPath = sec.Key("STATIC_ROOT_PATH").MustString(workDir)
	AppDataPath = sec.Key("APP_DATA_PATH").MustString("data")
	EnableGzip = sec.Key("ENABLE_GZIP").MustBool()

	sec = Cfg.Section("security")
	InstallLock = sec.Key("INSTALL_LOCK").MustBool()
	SecretKey = sec.Key("SECRET_KEY").String()
	LogInRememberDays = sec.Key("LOGIN_REMEMBER_DAYS").MustInt()
	CookieUserName = sec.Key("COOKIE_USERNAME").String()
	CookieRememberName = sec.Key("COOKIE_REMEMBER_NAME").String()
	CookieSecure = sec.Key("COOKIE_SECURE").MustBool(false)
	ReverseProxyAuthUser = sec.Key("REVERSE_PROXY_AUTHENTICATION_USER").MustString("X-WEBAUTH-USER")

	TimeFormat = map[string]string{
		"ANSIC":       time.ANSIC,
		"UnixDate":    time.UnixDate,
		"RubyDate":    time.RubyDate,
		"RFC822":      time.RFC822,
		"RFC822Z":     time.RFC822Z,
		"RFC850":      time.RFC850,
		"RFC1123":     time.RFC1123,
		"RFC1123Z":    time.RFC1123Z,
		"RFC3339":     time.RFC3339,
		"RFC3339Nano": time.RFC3339Nano,
		"Kitchen":     time.Kitchen,
		"Stamp":       time.Stamp,
		"StampMilli":  time.StampMilli,
		"StampMicro":  time.StampMicro,
		"StampNano":   time.StampNano,
	}[Cfg.Section("time").Key("FORMAT").MustString("RFC1123")]

	RunUser = Cfg.Section("").Key("RUN_USER").String()

	ProdMode = Cfg.Section("").Key("RUN_MODE").String() == "prod"

	HasRobotsTxt = com.IsFile(path.Join(CustomPath, "robots.txt"))

	initModulesConfig()
}

var Service struct {
	RegisterEmailConfirm   bool
	DisableRegistration    bool
	RequireSignInView      bool
	EnableCacheAvatar      bool
	EnableNotifyMail       bool
	EnableReverseProxyAuth bool
	ActiveCodeLives        int
	ResetPwdCodeLives      int
	EnableCaptcha          bool
}

func loadLogConfig() {

	hasConsole := false

	LogModes = strings.Split(Cfg.Section("log").Key("MODE").MustString("console"), ",")
	LogConfigs = make([]interface{}, len(LogModes))
	levelNames := map[string]log.LEVEL{
		"trace": log.TRACE,
		"info":  log.INFO,
		"warn":  log.WARN,
		"error": log.ERROR,
		"fatal": log.FATAL,
	}
	for i, mode := range LogModes {
		mode = strings.ToLower(strings.TrimSpace(mode))
		sec, err := Cfg.GetSection("log." + mode)
		if err != nil {
			log.Fatal(4, "Unknown logger mode: %s", mode)
		}

		name := Cfg.Section("log." + mode).Key("LEVEL").MustString("trace")
		level := levelNames[name]

		// Generate log configuration.
		switch log.MODE(mode) {
		case log.CONSOLE:
			hasConsole = true
			LogConfigs[i] = log.ConsoleConfig{
				Level:      level,
				BufferSize: Cfg.Section("log").Key("BUFFER_LEN").MustInt64(100),
			}

		case log.FILE:
			logPath := path.Join(LogRootPath, "goby.log")
			if err = os.MkdirAll(path.Dir(logPath), os.ModePerm); err != nil {
				log.Fatal(4, "Fail to create log directory '%s': %v", path.Dir(logPath), err)
			}

			LogConfigs[i] = log.FileConfig{
				Level:      level,
				BufferSize: Cfg.Section("log").Key("BUFFER_LEN").MustInt64(100),
				Filename:   logPath,
				FileRotationConfig: log.FileRotationConfig{
					Rotate:   sec.Key("LOG_ROTATE").MustBool(true),
					Daily:    sec.Key("DAILY_ROTATE").MustBool(true),
					MaxSize:  1 << uint(sec.Key("MAX_SIZE_SHIFT").MustInt(28)),
					MaxLines: sec.Key("MAX_LINES").MustInt64(1000000),
					MaxDays:  sec.Key("MAX_DAYS").MustInt64(7),
				},
			}

		case log.SLACK:
			LogConfigs[i] = log.SlackConfig{
				Level:      level,
				BufferSize: Cfg.Section("log").Key("BUFFER_LEN").MustInt64(100),
				URL:        sec.Key("URL").String(),
			}
		}

		log.New(log.MODE(mode), LogConfigs[i])
		log.Trace("Log Mode: %s (%s)", strings.Title(mode), strings.Title(name))
	}

	// Make sure everyone gets version info printed.
	log.Info("%s %s", AppName, AppVer)
	if !hasConsole {
		log.Delete(log.CONSOLE)
	}
}

func loadCacheConfig() {
	CacheAdapter = Cfg.Section("cache").Key("ADAPTER").In("memory", []string{"memory", "redis", "memcache"})
	CacheClientCheckUpdateTimeOut = Cfg.Section("cache").Key("CLIENT_CHECK_UPDATE_DATA_TIME_OUT").MustInt64(12 * 60 * 60 * 1000)
	CacheTokenTimeOut = Cfg.Section("cache").Key("TOKEN_DATA_TIME_OUT").MustInt64(60 * 60 * 1000)
	if EnableRedis {
		log.Info("Redis Supported")
	}
	if EnableMemcache {
		log.Info("Memcache Supported")
	}

	switch CacheAdapter {
	case "memory":
		CacheInternal = Cfg.Section("cache").Key("INTERVAL").MustInt(60)
	case "redis", "memcache":
		CacheConn = strings.Trim(Cfg.Section("cache").Key("HOST").String(), "\" ")
	default:
		log.Fatal(4, "Unknown cache adapter: %s", CacheAdapter)
	}

	log.Info("Cache Service Enabled")
}

func loadSessionConfig() {
	SessionConfig.Provider = Cfg.Section("session").Key("PROVIDER").In("memory",
		[]string{"memory", "file", "redis", "mysql"})
	SessionConfig.ProviderConfig = strings.Trim(Cfg.Section("session").Key("PROVIDER_CONFIG").String(), "\" ")
	SessionConfig.CookieName = Cfg.Section("session").Key("COOKIE_NAME").MustString("i_like_goby")
	SessionConfig.CookiePath = AppSubURL
	SessionConfig.Secure = Cfg.Section("session").Key("COOKIE_SECURE").MustBool()
	SessionConfig.Gclifetime = Cfg.Section("session").Key("GC_INTERVAL_TIME").MustInt64(86400)
	SessionConfig.Maxlifetime = Cfg.Section("session").Key("SESSION_LIFE_TIME").MustInt64(86400)

	log.Info("Session Service Enabled")
}

var DBConfig model.DbCfg

func loadDBConfigs() {
	sec := Cfg.Section("database")
	DBConfig.Type = sec.Key("DB_TYPE").String()
	switch DBConfig.Type {
	case "sqlite3":
		UseSQLite3 = true
	case "mysql":
		UseMySQL = true
	case "postgres":
		UsePostgreSQL = true
	case "tidb":
		UseTiDB = true
	}
	DBConfig.Host = sec.Key("HOST").String()
	DBConfig.Name = sec.Key("NAME").String()
	DBConfig.User = sec.Key("USER").String()
	if len(DBConfig.Passwd) == 0 {
		DBConfig.Passwd = sec.Key("PASSWD").String()
	}
	DBConfig.SSLMode = sec.Key("SSL_MODE").String()
	DBConfig.Path = sec.Key("PATH").MustString("data/goby.db")
	DBConfig.LogPath = LogRootPath
}

func loadServiceConfig() {
	sec := Cfg.Section("service")
	Service.ActiveCodeLives = sec.Key("ACTIVE_CODE_LIVE_MINUTES").MustInt(180)
	Service.ResetPwdCodeLives = sec.Key("RESET_PASSWD_CODE_LIVE_MINUTES").MustInt(180)
	Service.DisableRegistration = sec.Key("DISABLE_REGISTRATION").MustBool()
	Service.RequireSignInView = sec.Key("REQUIRE_SIGNIN_VIEW").MustBool()
	Service.EnableCacheAvatar = sec.Key("ENABLE_CACHE_AVATAR").MustBool()
	Service.EnableReverseProxyAuth = sec.Key("ENABLE_REVERSE_PROXY_AUTHENTICATION").MustBool()
	Service.EnableCaptcha = sec.Key("ENABLE_CAPTCHA").MustBool()

	Service.RegisterEmailConfirm = sec.Key("REGISTER_EMAIL_CONFIRM").MustBool()
	Service.EnableNotifyMail = sec.Key("REGISTER_EMAIL_CONFIRM").MustBool()
}

func loadAttachmentConfig() {
	sec := Cfg.Section("attachment")
	AttachmentAllowedTypes = strings.Replace(sec.Key("ALLOWED_TYPES").MustString("image/jpeg,image/png,application/zip"), "|", ",", -1)
	AttachmentMaxSize = sec.Key("MAX_SIZE").MustInt64(4)
	AttachmentMaxFiles = sec.Key("MAX_FILES").MustInt(5)
	AttachmentEnabled = sec.Key("ENABLE").MustBool(true)
	AttachmentPath = sec.Key("PATH").MustString(path.Join(AppDataPath, "attachments"))

	if err := os.MkdirAll(AttachmentPath, os.ModePerm); err != nil {
		log.Fatal(4, "Fail to create attachment directory '%s': %v", path.Dir(AttachmentPath), err)
	}
}

var Storage struct {
	StorageType   string
	StorageConfig *storage.Config
}

func loadStorageConfig() {
	sec := Cfg.Section("storage")
	Storage.StorageType = sec.Key("STORAGE_TYPE").In("local", []string{"local", "qiniu", "oss"})

	Storage.StorageConfig = &storage.Config{
		DownloadURL:      sec.Key("DOWNLOAD_URL").MustString(path.Join(AppURL, "download")),
		LocalStoragePath: sec.Key("STORAGE_PATH").MustString(path.Join(AppDataPath, "packages")),
		AccessKey:        sec.Key("ACCESS_KEY").MustString(""),
		SecretKey:        sec.Key("SECRET_KEY").MustString(""),
		Bucket:           sec.Key("BUCKET").MustString(""),
		Prefix:           sec.Key("PREFIX").MustString(""),
		QiNiuZone:        sec.Key("QN_ZONE").MustInt(1),
		OSSEndpoint:      sec.Key("OSS_ENDPOINT").MustString(""),
	}
}

var (
	MailService *mailer.Mailer
)

func loadMailConfig() {
	sec := Cfg.Section("mailer")
	// Check mailer setting.
	if !sec.Key("ENABLED").MustBool() {
		return
	}

	MailService = &mailer.Mailer{
		QueueLength:           sec.Key("SEND_BUFFER_LEN").MustInt(100),
		Name:                  sec.Key("NAME").MustString(AppName),
		Host:                  sec.Key("HOST").String(),
		User:                  sec.Key("USER").String(),
		Passwd:                sec.Key("PASSWD").String(),
		DisableHelo:           sec.Key("DISABLE_HELO").MustBool(),
		HeloHostname:          sec.Key("HELO_HOSTNAME").String(),
		SkipVerify:            sec.Key("SKIP_VERIFY").MustBool(),
		UseCertificate:        sec.Key("USE_CERTIFICATE").MustBool(),
		CertFile:              sec.Key("CERT_FILE").String(),
		KeyFile:               sec.Key("KEY_FILE").String(),
		EnableHTMLAlternative: sec.Key("ENABLE_HTML_ALTERNATIVE").MustBool(),
	}
	MailService.From = sec.Key("FROM").MustString(MailService.User)

	parsed, err := mail.ParseAddress(MailService.From)
	if err != nil {
		log.Fatal(4, "Invalid mailer.FROM (%s): %v", MailService.From, err)
	}
	MailService.FromEmail = parsed.Address

	log.Info("Mail Service Enabled")

}

var PackageConfig struct {
	EnableGoogleDiff bool
	MAXDiffCount     int
}

func loadPackageConfig() {
	sec := Cfg.Section("package")
	PackageConfig.EnableGoogleDiff = sec.Key("ENABLE_GOOLE_DIFF").MustBool(false)
	PackageConfig.MAXDiffCount = sec.Key("MAX_DIFF_COUNT").MustInt(1)
}

func initModulesConfig() {
	loadLogConfig()
	loadCacheConfig()
	loadSessionConfig()
	loadDBConfigs()
	loadStorageConfig()
	loadServiceConfig()
	loadAttachmentConfig()
	loadMailConfig()
}
