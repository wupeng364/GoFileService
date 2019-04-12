package configmodel

/**
 * 系统配置接口
 */
type cfInterface interface{
	
	InitConfig( ) error
	GetConfig(key string)string
	SetConfig(key string, value string)error
	
}

// 默认本地配置: {rundir}/config/config.json