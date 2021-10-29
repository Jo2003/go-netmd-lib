package netmd

import "github.com/enimatek-nl/gousb"

type Device struct {
	vendorId  gousb.ID
	deviceId  gousb.ID
	of_lp_enc bool
	name      string
}

var (
	Devices = [...]Device{
		{vendorId: 0x04dd, deviceId: 0x7202, of_lp_enc: false, name: "Sharp IM-MT899H"},
		{vendorId: 0x04dd, deviceId: 0x9013, of_lp_enc:  true, name: "Sharp IM-DR400/DR410/DR420"},
		{vendorId: 0x04dd, deviceId: 0x9014, of_lp_enc: false, name: "Sharp IM-DR80"},
		{vendorId: 0x054c, deviceId: 0x0034, of_lp_enc: false, name: "Sony PCLK-XX"},
		{vendorId: 0x054c, deviceId: 0x0036, of_lp_enc: false, name: "Sony"},
		{vendorId: 0x054c, deviceId: 0x0075, of_lp_enc: false, name: "Sony MZ-N1"},
		{vendorId: 0x054c, deviceId: 0x007c, of_lp_enc: false, name: "Sony"},
		{vendorId: 0x054c, deviceId: 0x0080, of_lp_enc: false, name: "Sony LAM-1"},
		{vendorId: 0x054c, deviceId: 0x0081, of_lp_enc:  true, name: "Sony MDS-JB980/JE780"},
		{vendorId: 0x054c, deviceId: 0x0084, of_lp_enc: false, name: "Sony MZ-N505"},
		{vendorId: 0x054c, deviceId: 0x0085, of_lp_enc: false, name: "Sony MZ-S1"},
		{vendorId: 0x054c, deviceId: 0x0086, of_lp_enc: false, name: "Sony MZ-N707"},
		{vendorId: 0x054c, deviceId: 0x008e, of_lp_enc: false, name: "Sony CMT-C7NT"},
		{vendorId: 0x054c, deviceId: 0x0097, of_lp_enc: false, name: "Sony PCGA-MDN1"},
		{vendorId: 0x054c, deviceId: 0x00ad, of_lp_enc: false, name: "Sony CMT-L7HD"},
		{vendorId: 0x054c, deviceId: 0x00c6, of_lp_enc: false, name: "Sony MZ-N10"},
		{vendorId: 0x054c, deviceId: 0x00c7, of_lp_enc: false, name: "Sony MZ-N910"},
		{vendorId: 0x054c, deviceId: 0x00c8, of_lp_enc: false, name: "Sony MZ-N710/NF810"},
		{vendorId: 0x054c, deviceId: 0x00c9, of_lp_enc: false, name: "Sony MZ-N510/N610"},
		{vendorId: 0x054c, deviceId: 0x00ca, of_lp_enc: false, name: "Sony MZ-NE410/NF520D"},
		{vendorId: 0x054c, deviceId: 0x00eb, of_lp_enc: false, name: "Sony MZ-NE810/NE910"},
		{vendorId: 0x054c, deviceId: 0x0101, of_lp_enc: false, name: "Sony LAM-10"},
		{vendorId: 0x054c, deviceId: 0x0113, of_lp_enc: false, name: "Aiwa AM-NX1"},
		{vendorId: 0x054c, deviceId: 0x013f, of_lp_enc: false, name: "Sony MDS-S500"},
		{vendorId: 0x054c, deviceId: 0x014c, of_lp_enc: false, name: "Aiwa AM-NX9"},
		{vendorId: 0x054c, deviceId: 0x017e, of_lp_enc: false, name: "Sony MZ-NH1"},
		{vendorId: 0x054c, deviceId: 0x0180, of_lp_enc: false, name: "Sony MZ-NH3D"},
		{vendorId: 0x054c, deviceId: 0x0182, of_lp_enc: false, name: "Sony MZ-NH900"},
		{vendorId: 0x054c, deviceId: 0x0184, of_lp_enc: false, name: "Sony MZ-NH700/NH800"},
		{vendorId: 0x054c, deviceId: 0x0186, of_lp_enc: false, name: "Sony MZ-NH600"},
		{vendorId: 0x054c, deviceId: 0x0187, of_lp_enc: false, name: "Sony MZ-NH600D"},
		{vendorId: 0x054c, deviceId: 0x0188, of_lp_enc: false, name: "Sony MZ-N920"},
		{vendorId: 0x054c, deviceId: 0x018a, of_lp_enc: false, name: "Sony LAM-3"},
		{vendorId: 0x054c, deviceId: 0x01e9, of_lp_enc: false, name: "Sony MZ-DH10P"},
		{vendorId: 0x054c, deviceId: 0x0219, of_lp_enc: false, name: "Sony MZ-RH10"},
		{vendorId: 0x054c, deviceId: 0x021b, of_lp_enc: false, name: "Sony MZ-RH710/MZ-RH910"},
		{vendorId: 0x054c, deviceId: 0x021d, of_lp_enc: false, name: "Sony CMT-AH10"},
		{vendorId: 0x054c, deviceId: 0x022c, of_lp_enc: false, name: "Sony CMT-AH10"},
		{vendorId: 0x054c, deviceId: 0x023c, of_lp_enc: false, name: "Sony DS-HMD1"},
		{vendorId: 0x054c, deviceId: 0x0286, of_lp_enc: false, name: "Sony MZ-RH1"},
	}
)
