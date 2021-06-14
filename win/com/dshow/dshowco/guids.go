package dshowco

import (
	"github.com/rodrigocfd/windigo/win/co"
)

// DirectShow COM CLSIDs.
const (
	CLSID_EnhancedVideoRenderer co.CLSID = "fa10746c-9b63-4b6c-bc49-fc300ea5f256"
	CLSID_FilterGraph           co.CLSID = "e436ebb3-524f-11ce-9f53-0020af0ba770"
	CLSID_MR_VideoRenderService co.CLSID = "1092a86c-ab1a-459a-a336-831fbc4d11ff"
	CLSID_VideoMixingRenderer9  co.CLSID = "51b4abf3-748f-4e3b-a276-c828330e926a"
)

// DirectShow COM IIDs.
const (
	IID_IBaseFilter            co.IID = "56a86895-0ad4-11ce-b03a-0020af0ba770"
	IID_IBasicAudio            co.IID = "56a868b3-0ad4-11ce-b03a-0020af0ba770"
	IID_IEnumFilters           co.IID = "56a86893-0ad4-11ce-b03a-0020af0ba770"
	IID_IEnumMediaTypes        co.IID = "89c31040-846b-11ce-97d3-00aa0055595a"
	IID_IFilterGraph           co.IID = "56a8689f-0ad4-11ce-b03a-0020af0ba770"
	IID_IGraphBuilder          co.IID = "56a868a9-0ad4-11ce-b03a-0020af0ba770"
	IID_IMediaControl          co.IID = "56a868b1-0ad4-11ce-b03a-0020af0ba770"
	IID_IMediaFilter           co.IID = "56a86899-0ad4-11ce-b03a-0020af0ba770"
	IID_IMediaSeeking          co.IID = "36b73880-c2c8-11cf-8b46-00805f6cef60"
	IID_IMFGetService          co.IID = "fa993888-4383-415a-a930-dd472a8cf6f7"
	IID_IMFVideoDisplayControl co.IID = "a490b1e4-ab84-4d31-a1b2-181e03b1077a"
	IID_IPersist               co.IID = "0000010c-0000-0000-c000-000000000046"
	IID_IPin                   co.IID = "56a86891-0ad4-11ce-b03a-0020af0ba770"
)
