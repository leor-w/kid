package kid

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/lsjhtang/kid/logger"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/database/mysql"
	"github.com/spf13/cobra"
)

var (
	customConf bool
	database   string
	address    string
	port       int
	user       string
	password   string
)

const template = `{{if .HasAvailableSubCommands}}可用命令:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

全局 Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

其他帮助信息:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

使用 "{{.CommandPath}} [command] --help" 获取有关命令的更多信息.{{end}}
`

// isTestEnvironment 检测当前是否处于 go test 环境中
func isTestEnvironment() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}

// CommandLine 命令行处理
func (kid *Kid) CommandLine() error {
	// 检测是否在测试环境中
	if isTestEnvironment() {
		return nil // 直接返回，跳过命令行执行
	}

	var (
		hasHelp bool
	)
	rootCmd := &cobra.Command{Use: kid.Options.Name}
	cobra.EnableCommandSorting = true
	cobra.EnablePrefixMatching = true
	rootCmd.AddCommand(kid.NewDbCommand())
	rootCmd.AddCommand(kid.NewServerCommand())
	rootCmd.SetHelpCommand(kid.NewHelpCommand())
	helpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		hasHelp = true
		helpFunc(cmd, args)
	})
	rootCmd.SetHelpTemplate(template)
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	if hasHelp {
		syscall.Exit(0)
	}
	return nil
}

// NewDbCommand 数据库相关命令
func (kid *Kid) NewDbCommand() *cobra.Command {
	dbCmd := &cobra.Command{
		Use:   "database",
		Short: "数据库操作命令",
		Long:  "可以进行数据库配置及数据表迁移等功能，可以使用 -c 指定配置文件，也可以使用选项指定连接参数。",
		Run: func(cmd *cobra.Command, args []string) {
			migrate, err := cmd.Flags().GetBool("automigrate")
			if err != nil {
				errorPrint(fmt.Sprintf("获取是否同步数据库选项错误: %s", err.Error()))
				syscall.Exit(1)
			}
			kid.Options.AutoMigrate = migrate
			conf, err := cmd.Flags().GetString("config.yaml")
			if err != nil {
				panic(fmt.Sprintf("读取配置文件错误: %s", err.Error()))
			}
			closeFKCheck, err := cmd.Flags().GetBool("closefkcheck")
			if err != nil {
				panic(fmt.Sprintf("获取是否关闭外键检查错误: %s", err.Error()))
			}
			if len(conf) > 0 {
				customConf = true
				kid.Options.Configs = []string{conf}
			} else {
				database, err = cmd.Flags().GetString("database")
				if err != nil {
					errorPrint(fmt.Sprintf("解析数据库名称错误: %s", err.Error()))
					syscall.Exit(1)
				}
				address, err = cmd.Flags().GetString("address")
				if err != nil {
					errorPrint(fmt.Sprintf("解析数据库连接地址错误: %s", err.Error()))
					syscall.Exit(1)
				}
				port, err = cmd.Flags().GetInt("port")
				if err != nil {
					errorPrint(fmt.Sprintf("解析数据库连接端口错误: %s", err.Error()))
					syscall.Exit(1)
				}
				user, err = cmd.Flags().GetString("user")
				if err != nil {
					errorPrint(fmt.Sprintf("解析数据库连接用户错误: %s", err.Error()))
					syscall.Exit(1)
				}
				password, err = cmd.Flags().GetString("password")
				if err != nil {
					errorPrint(fmt.Sprintf("解析数据库连接用户密码错误: %s", err.Error()))
					syscall.Exit(1)
				}
			}
			if migrate {
				if err := kid.loadConfig(); err != nil {
					fmt.Println(fmt.Sprintf("读取配置文件错误: %s", err.Error()))
					syscall.Exit(1)
				}
				if err := kid.loadLogger(); err != nil {
					fmt.Println(fmt.Sprintf("初始化日志错误: %s", err.Error()))
					syscall.Exit(1)
				}
				kid.AutoMigrate(closeFKCheck)
			}
			syscall.Exit(0)
		},
		SilenceErrors: true,
	}
	dbCmd.Flags().BoolP("automigrate", "m", false, "进行数据库迁移操作, 此操作将会同步项目中注册的 models 到数据库表")
	dbCmd.Flags().BoolP("closefkcheck", "f", false, "是否关闭外键检查")
	dbCmd.Flags().StringP("config.yaml", "c", "", "指定配置文件")
	dbCmd.Flags().StringP("database", "n", "", "需要同步的库")
	dbCmd.Flags().StringP("address", "a", "", "数据库连接地址")
	dbCmd.Flags().IntP("port", "i", 3306, "数据库连接端口")
	dbCmd.Flags().StringP("user", "u", "", "连接用户")
	dbCmd.Flags().StringP("password", "p", "", "连接用户密码")
	return dbCmd
}

// NewServerCommand 服务启动相关命令
func (kid *Kid) NewServerCommand() *cobra.Command {
	srvCmd := &cobra.Command{
		Use:   "server",
		Short: "服务启动相关命令",
		Long:  "使用该指令可以指定启动的配置文件等选项",
		Run: func(cmd *cobra.Command, args []string) {
			confs, err := cmd.Flags().GetStringSlice("config.yaml")
			if err != nil {
				panic(fmt.Sprintf("服务启动错误: 解析手动指定配置文件错误: %s", err.Error()))
			}
			if len(confs) > 0 {
				kid.Options.Configs = confs
			}
		},
		SilenceErrors: true,
	}
	srvCmd.Flags().StringSliceP("config.yaml", "c", []string{}, "服务启动的 yaml 配置文件，多个配置文件可以使用 ',' 分隔(默认为 ./conf/config.yaml.yaml)")
	return srvCmd
}

func (kid *Kid) NewHelpCommand() *cobra.Command {
	helpCmd := &cobra.Command{
		Use:   "help",
		Short: "关于任何命令的帮助",
		Long:  "关于任何命令的帮助",
	}
	return helpCmd
}

// AutoMigrate 自动迁移数据库
func (kid *Kid) AutoMigrate(closeFKCheck bool) {
	var db *mysql.MySQL
	if customConf {
		conf := config.New(config.WithProviders(kid.Options.Configs))
		db = mysql.New(
			mysql.WithHost(conf.GetString("mysql.host")),
			mysql.WithPort(conf.GetInt("mysql.port")),
			mysql.WithUser(conf.GetString("mysql.user")),
			mysql.WithPassword(conf.GetString("mysql.password")),
			mysql.WithDb(conf.GetString("mysql.database")),
			mysql.WithLogLevel(4),
			mysql.WithCloseFKCheck(closeFKCheck),
		)
	} else {
		db = mysql.New(
			mysql.WithHost(address),
			mysql.WithPort(port),
			mysql.WithUser(user),
			mysql.WithPassword(password),
			mysql.WithDb(database),
			mysql.WithLogLevel(4),
			mysql.WithCloseFKCheck(closeFKCheck),
		)
	}
	if kid.Options.DatabaseMigrate == nil {
		logger.Errorf("没有配置有效的模型提供者")
		syscall.Exit(1)
	}
	printInfo("==============================数据库同步开始==============================")
	if err := db.AutoMigrate(kid.Options.DatabaseMigrate.Models()...); err != nil {
		errorPrint(fmt.Sprintf("数据库表同步错误: %s", err.Error()))
		syscall.Exit(1)
	}
	printInfo("==============================数据库同步完成==============================")
	syscall.Exit(0)
}

func printInfo(args ...interface{}) {
	fmt.Println(args...)
}

func errorPrint(args ...interface{}) {
	args = append([]interface{}{"Error: "}, args...)
	fmt.Println(args...)
}
