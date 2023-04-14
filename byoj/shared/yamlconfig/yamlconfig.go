package yamlconfig

import (
	"byoj/model"
	"byoj/shared/server"
	"byoj/utils/logs"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"go.uber.org/zap"
)

type Configuration struct {
	Server   server.Server     `yaml:"server"`
	Database model.Database `yaml:"database"`
}

func ConfigLoad(path string) (Configuration, error) {
	// 设置选项支持 ENV 解析
	config.WithOptions(config.ParseEnv)

	// 添加驱动程序以支持 yaml 内容解析
	config.AddDriver(yaml.Driver)
	config.WithOptions(func(opt *config.Options) {
		opt.DecoderConfig.TagName = "yaml"
	})

	configuration := Configuration{}
	err := config.LoadFiles(path)
	if err != nil {
		logs.Error("Read config file from "+path+"failed. ", zap.Error(err))
		return configuration, err
	}

	err = config.Decode(&configuration)
	if err != nil {
		logs.Error("Decode config file from "+path+"failed. ", zap.Error(err))
		return configuration, err
	}
	logs.Info("Read config file from " + path)

	return configuration, nil
}
