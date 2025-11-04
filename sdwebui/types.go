package sdwebui

type TextToImageRequest struct {
	Prompt              string                 `json:"prompt" jsonschema:"提示词,描述要生成的图片内容"`
	NegativePrompt      string                 `json:"negative_prompt,omitempty" jsonschema:"负面提示词,不希望出现在图片中的内容"`
	Width               int                    `json:"width,omitempty" jsonschema:"图片宽度,生成图片的宽度（像素）"`
	Height              int                    `json:"height,omitempty" jsonschema:"图片高度,生成图片的高度（像素）"`
	Steps               int                    `json:"steps,omitempty" jsonschema:"采样步数,扩散过程的迭代次数"`
	SamplerName         string                 `json:"sampler_name,omitempty" jsonschema:"采样器名称,使用的采样算法"`
	Seed                int64                  `json:"seed,omitempty" jsonschema:"随机种子,控制生成结果的随机性"`
	CFGScale            float64                `json:"cfg_scale,omitempty" jsonschema:"提示词相关性,控制提示词对生成结果的影响程度"`
	BatchSize           int                    `json:"batch_size,omitempty" jsonschema:"批次大小,单次生成的图片数量"`
	NIter               int                    `json:"n_iter,omitempty" jsonschema:"批次数量,生成批次的次数"`
	EnableHR            bool                   `json:"enable_hr,omitempty" jsonschema:"是否启用高分辨率修复,是否启用高分辨率放大"`
	HRScale             float64                `json:"hr_scale,omitempty" jsonschema:"高分辨率修复比例,高分辨率放大的比例"`
	HRSamplerName       string                 `json:"hr_sampler_name,omitempty" jsonschema:"高分辨率修复采样器,高分辨率阶段使用的采样器"`
	HRSteps             int                    `json:"hr_steps,omitempty" jsonschema:"高分辨率修复步数,高分辨率阶段的采样步数"`
	HRDenoisingStrength float64                `json:"hr_denoising_strength,omitempty" jsonschema:"高分辨率修复去噪强度,高分辨率阶段的去噪强度"`
	HRUpscaler          string                 `json:"hr_upscaler,omitempty" jsonschema:"高分辨率修复放大算法,高分辨率放大使用的算法"`
	RestoreFaces        bool                   `json:"restore_faces,omitempty" jsonschema:"是否使用面部修复,是否启用面部修复功能"`
	Tiling              bool                   `json:"tiling,omitempty" jsonschema:"是否使用平铺,是否生成可平铺的图片"`
	OverrideSettings    map[string]interface{} `json:"override_settings,omitempty" jsonschema:"是否覆盖设置,覆盖默认设置的自定义参数"`
	ScriptArgs          []interface{}          `json:"script_args,omitempty" jsonschema:"脚本参数,脚本功能的参数列表"`
	ScriptName          string                 `json:"script_name,omitempty" jsonschema:"脚本名称,要使用的脚本名称"`

	// ControlNet 相关参数
	ControlNetEnabled bool             `json:"controlnet_enabled,omitempty" jsonschema:"是否启用ControlNet,是否启用ControlNet扩展"`
	ControlNetUnits   []ControlNetUnit `json:"controlnet_units,omitempty" jsonschema:"ControlNet单元列表,一个或多个ControlNet配置单元"`
}

type TextToImageResponse struct {
	Images     []string               `json:"images" jsonschema:"生成的图片列表,生成的图片列表（图片url）"`
	Parameters map[string]interface{} `json:"parameters" jsonschema:"生成参数,生成参数信息"`
	Info       string                 `json:"info" jsonschema:"生成信息,详细的生成信息"`
}

// ControlNetUnit 定义了单个 ControlNet 的配置
type ControlNetUnit struct {
	// 输入图像，通常是base64（不含前缀）或可访问的URL
	InputImage string `json:"input_image,omitempty" jsonschema:"输入图像,作为ControlNet的条件图像（base64或URL）"`
	// 可选遮罩
	Mask string `json:"mask,omitempty" jsonschema:"遮罩,可选的遮罩图像（base64或URL）"`
	// 预处理模块（如: canny, depth, softedge 等）
	Module string `json:"module,omitempty" jsonschema:"预处理模块,如canny/depth/softedge等"`
	// ControlNet 模型名称
	Model string `json:"model,omitempty" jsonschema:"ControlNet模型,ControlNet模型名称"`
	// 影响权重（0~2）
	Weight float64 `json:"weight,omitempty" jsonschema:"权重,ControlNet影响强度(0-2)"`
	// 尺寸调整模式（如: Resize, Crop, Envelope）
	ResizeMode string `json:"resize_mode,omitempty" jsonschema:"缩放模式,输入图像到目标尺寸的处理方式"`
	// 是否低显存模式
	LowVram bool `json:"lowvram,omitempty" jsonschema:"低显存,是否启用低显存模式"`
	// 预处理分辨率（像素）
	ProcessorRes int `json:"processor_res,omitempty" jsonschema:"预处理分辨率,预处理器使用的分辨率"`
	// A/B 阈值，随不同模块含义不同
	ThresholdA float64 `json:"threshold_a,omitempty" jsonschema:"阈值A,模块相关阈值A"`
	ThresholdB float64 `json:"threshold_b,omitempty" jsonschema:"阈值B,模块相关阈值B"`
	// 引导强度与起止步（0.0-1.0）
	Guidance      float64 `json:"guidance,omitempty" jsonschema:"引导强度,ControlNet引导强度"`
	GuidanceStart float64 `json:"guidance_start,omitempty" jsonschema:"引导起始,在扩散过程中的起始比率(0-1)"`
	GuidanceEnd   float64 `json:"guidance_end,omitempty" jsonschema:"引导结束,在扩散过程中的结束比率(0-1)"`
	// 控制模式（如: Balanced, My prompt is more important, ControlNet is more important）
	ControlMode string `json:"control_mode,omitempty" jsonschema:"控制模式,提示词与ControlNet的相对重要性"`
	// 像素精确（尽量保留输入边缘/结构）
	PixelPerfect bool `json:"pixel_perfect,omitempty" jsonschema:"像素精确,是否启用像素精确模式"`
	// 可选：多图输入
	InputImages []string `json:"input_images,omitempty" jsonschema:"多图输入,可选的多张条件图像列表"`
}

type SdModelsResponse struct {
	Models []SdModel `json:"models" jsonschema:"模型列表,模型列表"`
}

type SdModel struct {
	Title             string   `json:"title" jsonschema:"模型标题,模型标题"`
	ModelName         string   `json:"model_name" jsonschema:"模型名称,模型名称"`
	Hash              string   `json:"hash" jsonschema:"模型哈希值,模型哈希值"`
	Filename          string   `json:"filename" jsonschema:"模型文件名,模型文件名"`
	Config            string   `json:"config" jsonschema:"配置文件,配置文件"`
	Type              string   `json:"type" jsonschema:"模型类型,模型类型"`
	Size              int64    `json:"size" jsonschema:"模型大小,模型大小（字节）"`
	Active            bool     `json:"active" jsonschema:"是否激活,是否激活"`
	Thumbnail         string   `json:"thumbnail,omitempty" jsonschema:"缩略图URL,缩略图URL"`
	Description       string   `json:"description,omitempty" jsonschema:"模型描述,模型描述"`
	Tags              []string `json:"tags,omitempty" jsonschema:"模型标签,模型标签"`
	SupportedSamplers []string `json:"supported_samplers,omitempty" jsonschema:"支持的采样器,支持的采样器"`
	SupportedSizes    []string `json:"supported_sizes,omitempty" jsonschema:"支持的尺寸,支持的尺寸"`
	Version           string   `json:"version,omitempty" jsonschema:"模型版本,模型版本"`
	TrainingData      string   `json:"training_data,omitempty" jsonschema:"训练数据,训练数据"`
	TrainingSteps     int64    `json:"training_steps,omitempty" jsonschema:"训练步数,训练步数"`
	BaseModel         string   `json:"base_model,omitempty" jsonschema:"基础模型,基础模型"`
}

type SwitchModelRequest struct {
	SdModelCheckpoint string `json:"sd_model_checkpoint" jsonschema:"模型名称,模型名称"`
}

type SwitchModelResponse struct {
	Success bool   `json:"success" jsonschema:"是否成功,是否成功切换模型"`
	Message string `json:"message,omitempty" jsonschema:"消息,操作结果消息"`
}
