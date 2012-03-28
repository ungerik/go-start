// Internationalization support - not ready yet.
package i18n

// This can be used as independent library

var iso3166_1_alpha2 map[string]string

func EnglishCountryName(code string) string {
	name, ok := Countries()[code]
	if !ok {
		return code
	}
	return name
}

// Countries returns a map of ISO 3166-1 alpha-2 country codes
// to the corresponding english country name
func Countries() map[string]string {
	if iso3166_1_alpha2 == nil {
		iso3166_1_alpha2 = map[string]string{
			"AD": "Andorra",
			"AE": "United Arab Emirates",
			"AF": "Afghanistan",
			"AG": "Antigua and Barbuda",
			"AI": "Anguilla",
			"AL": "Albania",	
			"AM": "Armenia",	
			"AO": "Angola",	
			"AQ": "Antarctica",
			"AR": "Argentina",	
			"AS": "American Samoa",	
			"AT": "Austria",	
			"AU": "Australia",
			"AW": "Aruba",	
			"AX": "Åland Islands",	
			"AZ": "Azerbaijan",	
			"BA": "Bosnia and Herzegovina",	
			"BB": "Barbados",	
			"BD": "Bangladesh",	
			"BE": "Belgium",	
			"BF": "Burkina Faso",
			"BG": "Bulgaria",	
			"BH": "Bahrain",	
			"BI": "Burundi",	
			"BJ": "Benin",
			"BL": "Saint Barthélemy",	
			"BM": "Bermuda",	
			"BN": "Brunei Darussalam",
			"BO": "Bolivia",
			"BQ": "Bonaire, Sint Eustatius and Saba",
			"BR": "Brazil",	
			"BS": "Bahamas",	
			"BT": "Bhutan",	
			"BV": "Bouvet Island",	
			"BW": "Botswana",	
			"BY": "Belarus",
			"BZ": "Belize",	
			"CA": "Canada",	
			"CC": "Cocos (Keeling) Islands",	
			"CD": "Congo, the Democratic Republic of the",
			"CF": "Central African Republic",	
			"CG": "Congo",	
			"CH": "Switzerland",
			"CI": "Côte d'Ivoire",	
			"CK": "Cook Islands",	
			"CL": "Chile",	
			"CM": "Cameroon",	
			"CN": "China",	
			"CO": "Colombia",	
			"CR": "Costa Rica",	
			"CU": "Cuba",	
			"CV": "Cape Verde",	
			"CW": "Curaçao",	
			"CX": "Christmas Island",	
			"CY": "Cyprus",	
			"CZ": "Czech Republic",	
			"DE": "Germany",
			"DJ": "Djibouti",
			"DK": "Denmark",	
			"DM": "Dominica",	
			"DO": "Dominican Republic",
			/*	
			"DZ": "Algeria	1974	.dz	ISO 3166-2:DZ	Code taken from name in Kabyle: Dzayer
			"EC": "Ecuador	1974	.ec	ISO 3166-2:EC	
			"EE": "Estonia	1992	.ee	ISO 3166-2:EE	Code taken from name in Estonian: Eesti
			"EG": "Egypt	1974	.eg	ISO 3166-2:EG	
			"EH": "Western Sahara	1974	.eh	ISO 3166-2:EH	Previous ISO country name: Spanish Sahara (code taken from name in Spanish: Sahara español)
			"ER": "Eritrea	1993	.er	ISO 3166-2:ER	
			"ES": "Spain	1974	.es	ISO 3166-2:ES	Code taken from name in Spanish: España
			"ET": "Ethiopia	1974	.et	ISO 3166-2:ET	
			"FI": "Finland	1974	.fi	ISO 3166-2:FI	
			"FJ": "Fiji	1974	.fj	ISO 3166-2:FJ	
			"FK": "Falkland Islands (Malvinas)	1974	.fk	ISO 3166-2:FK	
			"FM": "Micronesia, Federated States of	1986	.fm	ISO 3166-2:FM	Previous ISO country name: Micronesia
			"FO": "Faroe Islands	1974	.fo	ISO 3166-2:FO	
			"FR": "France	1974	.fr	ISO 3166-2:FR	Includes Clipperton Island
			"GA": "Gabon	1974	.ga	ISO 3166-2:GA	
			"GB": "United Kingdom	1974	.gb
			"GD": "Grenada	1974	.gd	ISO 3166-2:GD	
			"GE": "Georgia	1992	.ge	ISO 3166-2:GE	GE previously represented Gilbert and Ellice Islands
			"GF": "French Guiana	1974	.gf	ISO 3166-2:GF	Code taken from name in French: Guyane française
			"GG": "Guernsey	2006	.gg	ISO 3166-2:GG	
			"GH": "Ghana	1974	.gh	ISO 3166-2:GH	
			"GI": "Gibraltar	1974	.gi	ISO 3166-2:GI	
			"GL": "Greenland	1974	.gl	ISO 3166-2:GL	
			"GM": "Gambia	1974	.gm	ISO 3166-2:GM	
			"GN": "Guinea	1974	.gn	ISO 3166-2:GN	
			"GP": "Guadeloupe	1974	.gp	ISO 3166-2:GP	
			"GQ": "Equatorial Guinea	1974	.gq	ISO 3166-2:GQ	Code taken from name in French: Guinée équatoriale
			"GR": "Greece	1974	.gr	ISO 3166-2:GR	
			"GS": "South Georgia and the South Sandwich Islands	1993	.gs	ISO 3166-2:GS	
			"GT": "Guatemala	1974	.gt	ISO 3166-2:GT	
			"GU": "Guam	1974	.gu	ISO 3166-2:GU	
			"GW": "Guinea-Bissau	1974	.gw	ISO 3166-2:GW	
			"GY": "Guyana	1974	.gy	ISO 3166-2:GY	
			"HK": "Hong Kong	1974	.hk	ISO 3166-2:HK	
			"HM": "Heard Island and McDonald Islands	1974	.hm	ISO 3166-2:HM	
			"HN": "Honduras	1974	.hn	ISO 3166-2:HN	
			"HR": "Croatia	1992	.hr	ISO 3166-2:HR	Code taken from name in Croatian: Hrvatska
			"HT": "Haiti	1974	.ht	ISO 3166-2:HT	
			"HU": "Hungary	1974	.hu	ISO 3166-2:HU	
			"ID": "Indonesia	1974	.id	ISO 3166-2:ID	
			"IE": "Ireland	1974	.ie	ISO 3166-2:IE	
			"IL": "Israel	1974	.il	ISO 3166-2:IL	
			"IM": "Isle of Man	2006	.im	ISO 3166-2:IM	
			"IN": "India	1974	.in	ISO 3166-2:IN	
			"IO": "British Indian Ocean Territory	1974	.io	ISO 3166-2:IO	
			"IQ": "Iraq	1974	.iq	ISO 3166-2:IQ	
			"IR": "Iran, Islamic Republic of	1974	.ir	ISO 3166-2:IR	ISO country name follows UN designation (common name: Iran)
			"IS": "Iceland	1974	.is	ISO 3166-2:IS	Code taken from name in Icelandic: Ísland
			"IT": "Italy	1974	.it	ISO 3166-2:IT	
			"JE": "Jersey	2006	.je	ISO 3166-2:JE	
			"JM": "Jamaica	1974	.jm	ISO 3166-2:JM	
			"JO": "Jordan	1974	.jo	ISO 3166-2:JO	
			"JP": "Japan	1974	.jp	ISO 3166-2:JP	
			"KE": "Kenya	1974	.ke	ISO 3166-2:KE	
			"KG": "Kyrgyzstan	1992	.kg	ISO 3166-2:KG	
			"KH": "Cambodia	1974	.kh	ISO 3166-2:KH	Code taken from former name: Khmer Republic
			"KI": "Kiribati	1979	.ki	ISO 3166-2:KI	
			"KM": "Comoros	1974	.km	ISO 3166-2:KM	Code taken from name in Comorian: Komori
			"KN": "Saint Kitts and Nevis	1974	.kn	ISO 3166-2:KN	Previous ISO country name: Saint Kitts-Nevis-Anguilla
			"KP": "Korea, Democratic People's Republic of	1974	.kp	ISO 3166-2:KP	ISO country name follows UN designation (common name: North Korea)
			"KR": "Korea, Republic of	1974	.kr	ISO 3166-2:KR	ISO country name follows UN designation (common name: South Korea)
			"KW": "Kuwait	1974	.kw	ISO 3166-2:KW	
			"KY": "Cayman Islands	1974	.ky	ISO 3166-2:KY	
			"KZ": "Kazakhstan	1992	.kz	ISO 3166-2:KZ	Previous ISO country name: Kazakstan
			"LA": "Lao People's Democratic Republic	1974	.la	ISO 3166-2:LA	ISO country name follows UN designation (common name: Laos)
			"LB": "Lebanon	1974	.lb	ISO 3166-2:LB	
			"LC": "Saint Lucia	1974	.lc	ISO 3166-2:LC	
			"LI": "Liechtenstein	1974	.li	ISO 3166-2:LI	
			"LK": "Sri Lanka	1974	.lk	ISO 3166-2:LK	
			"LR": "Liberia	1974	.lr	ISO 3166-2:LR	
			"LS": "Lesotho	1974	.ls	ISO 3166-2:LS	
			"LT": "Lithuania	1992	.lt	ISO 3166-2:LT	
			"LU": "Luxembourg	1974	.lu	ISO 3166-2:LU	
			"LV": "Latvia	1992	.lv	ISO 3166-2:LV	
			"LY": "Libyan Arab Jamahiriya	1974	.ly	ISO 3166-2:LY	ISO country name follows UN designation (common name: Libya)
			"MA": "Morocco	1974	.ma	ISO 3166-2:MA	Code taken from name in French: Maroc
			"MC": "Monaco	1974	.mc	ISO 3166-2:MC	
			"MD": "Moldova, Republic of	1992	.md	ISO 3166-2:MD	ISO country name follows UN designation (common name and previous ISO country name: Moldova)
			"ME": "Montenegro	2006	.me	ISO 3166-2:ME	
			"MF": "Saint Martin (French part)	2007	.mf	ISO 3166-2:MF	The Dutch part of Saint Martin island is assigned code SX
			"MG": "Madagascar	1974	.mg	ISO 3166-2:MG	
			"MH": "Marshall Islands	1986	.mh	ISO 3166-2:MH	
			"MK": "Macedonia, the former Yugoslav Republic of	1993	.mk	ISO 3166-2:MK	ISO country name follows UN designation (due to Macedonia naming dispute; official name used by country itself: Republic of Macedonia)
			"ML": "Mali	1974	.ml	ISO 3166-2:ML	
			"MM": "Myanmar	1989	.mm	ISO 3166-2:MM	Name changed from Burma (BU)
			"MN": "Mongolia	1974	.mn	ISO 3166-2:MN	
			"MO": "Macao	1974	.mo	ISO 3166-2:MO	Previous ISO country name: Macau
			"MP": "Northern Mariana Islands	1986	.mp	ISO 3166-2:MP	
			"MQ": "Martinique	1974	.mq	ISO 3166-2:MQ	
			"MR": "Mauritania	1974	.mr	ISO 3166-2:MR	
			"MS": "Montserrat	1974	.ms	ISO 3166-2:MS	
			"MT": "Malta	1974	.mt	ISO 3166-2:MT	
			"MU": "Mauritius	1974	.mu	ISO 3166-2:MU	
			"MV": "Maldives	1974	.mv	ISO 3166-2:MV	
			"MW": "Malawi	1974	.mw	ISO 3166-2:MW	
			"MX": "Mexico	1974	.mx	ISO 3166-2:MX	
			"MY": "Malaysia	1974	.my	ISO 3166-2:MY	
			"MZ": "Mozambique	1974	.mz	ISO 3166-2:MZ	
			"NA": "Namibia	1974	.na	ISO 3166-2:NA	
			"NC": "New Caledonia	1974	.nc	ISO 3166-2:NC	
			"NE": "Niger	1974	.ne	ISO 3166-2:NE	
			"NF": "Norfolk Island	1974	.nf	ISO 3166-2:NF	
			"NG": "Nigeria	1974	.ng	ISO 3166-2:NG	
			"NI": "Nicaragua	1974	.ni	ISO 3166-2:NI	
			"NL": "Netherlands	1974	.nl	ISO 3166-2:NL	
			"NO": "Norway	1974	.no	ISO 3166-2:NO	
			"NP": "Nepal	1974	.np	ISO 3166-2:NP	
			"NR": "Nauru	1974	.nr	ISO 3166-2:NR	
			"NU": "Niue	1974	.nu	ISO 3166-2:NU	
			"NZ": "New Zealand	1974	.nz	ISO 3166-2:NZ	
			"OM": "Oman	1974	.om	ISO 3166-2:OM	
			"PA": "Panama	1974	.pa	ISO 3166-2:PA	
			"PE": "Peru	1974	.pe	ISO 3166-2:PE	
			"PF": "French Polynesia	1974	.pf	ISO 3166-2:PF	Code taken from name in French: Polynésie française
			"PG": "Papua New Guinea	1974	.pg	ISO 3166-2:PG	
			"PH": "Philippines	1974	.ph	ISO 3166-2:PH	
			"PK": "Pakistan	1974	.pk	ISO 3166-2:PK	
			"PL": "Poland	1974	.pl	ISO 3166-2:PL	
			"PM": "Saint Pierre and Miquelon	1974	.pm	ISO 3166-2:PM	
			"PN": "Pitcairn	1974	.pn	ISO 3166-2:PN	
			"PR": "Puerto Rico	1974	.pr	ISO 3166-2:PR	
			"PS": "Palestinian Territory, Occupied	1999	.ps	ISO 3166-2:PS	Consists of the West Bank and the Gaza Strip
			"PT": "Portugal	1974	.pt	ISO 3166-2:PT	
			"PW": "Palau	1986	.pw	ISO 3166-2:PW	
			"PY": "Paraguay	1974	.py	ISO 3166-2:PY	
			"QA": "Qatar	1974	.qa	ISO 3166-2:QA	
			"RE": "Réunion	1974	.re	ISO 3166-2:RE	
			"RO": "Romania	1974	.ro	ISO 3166-2:RO	
			"RS": "Serbia	2006	.rs	ISO 3166-2:RS	Code taken from official name: Republic of Serbia (see Serbian country codes)
			"RU": "Russian Federation	1992	.ru	ISO 3166-2:RU	ISO country name follows UN designation (common name: Russia)
			"RW": "Rwanda	1974	.rw	ISO 3166-2:RW	
			"SA": "Saudi Arabia	1974	.sa	ISO 3166-2:SA	
			"SB": "Solomon Islands	1974	.sb	ISO 3166-2:SB	Code taken from former name: British Solomon Islands
			"SC": "Seychelles	1974	.sc	ISO 3166-2:SC	
			"SD": "Sudan	1974	.sd	ISO 3166-2:SD	
			"SE": "Sweden	1974	.se	ISO 3166-2:SE	
			"SG": "Singapore	1974	.sg	ISO 3166-2:SG	
			"SH": "Saint Helena, Ascension and Tristan da Cunha	1974	.sh	ISO 3166-2:SH	Previous ISO country name: Saint Helena
			"SI": "Slovenia	1992	.si	ISO 3166-2:SI	
			"SJ": "Svalbard and Jan Mayen	1974	.sj	ISO 3166-2:SJ	Consists of two arctic territories of Norway: Svalbard and Jan Mayen
			"SK": "Slovakia	1993	.sk	ISO 3166-2:SK	SK previously represented Sikkim
			"SL": "Sierra Leone	1974	.sl	ISO 3166-2:SL	
			"SM": "San Marino	1974	.sm	ISO 3166-2:SM	
			"SN": "Senegal	1974	.sn	ISO 3166-2:SN	
			"SO": "Somalia	1974	.so	ISO 3166-2:SO	
			"SR": "Suriname	1974	.sr	ISO 3166-2:SR	
			"SS": "South Sudan	2011	.ss	ISO 3166-2:SS	
			"ST": "Sao Tome and Principe	1974	.st	ISO 3166-2:ST	
			"SV": "El Salvador	1974	.sv	ISO 3166-2:SV	
			"SX": "Sint Maarten (Dutch part)	2010	.sx	ISO 3166-2:SX	The French part of Saint Martin island is assigned code MF
			"SY": "Syrian Arab Republic	1974	.sy	ISO 3166-2:SY	ISO country name follows UN designation (common name: Syria)
			"SZ": "Swaziland	1974	.sz	ISO 3166-2:SZ	
			"TC": "Turks and Caicos Islands	1974	.tc	ISO 3166-2:TC	
			"TD": "Chad	1974	.td	ISO 3166-2:TD	Code taken from name in French: Tchad
			"TF": "French Southern Territories	1979	.tf	ISO 3166-2:TF	Covers the French Southern and Antarctic Lands except Adélie Land
			"TG": "Togo	1974	.tg	ISO 3166-2:TG	
			"TH": "Thailand	1974	.th	ISO 3166-2:TH	
			"TJ": "Tajikistan	1992	.tj	ISO 3166-2:TJ	
			"TK": "Tokelau	1974	.tk	ISO 3166-2:TK	
			"TL": "Timor-Leste	2002	.tl	ISO 3166-2:TL	Name changed from East Timor (TP)
			"TM": "Turkmenistan	1992	.tm	ISO 3166-2:TM	
			"TN": "Tunisia	1974	.tn	ISO 3166-2:TN	
			"TO": "Tonga	1974	.to	ISO 3166-2:TO	
			"TR": "Turkey	1974	.tr	ISO 3166-2:TR	
			"TT": "Trinidad and Tobago	1974	.tt	ISO 3166-2:TT	
			"TV": "Tuvalu	1979	.tv	ISO 3166-2:TV	
			"TW": "Taiwan, Province of China	1974	.tw	ISO 3166-2:TW	Covers the current jurisdiction of the Republic of China except Kinmen and Lienchiang
			"TZ": "Tanzania, United Republic of	1974	.tz	ISO 3166-2:TZ	ISO country name follows UN designation (common name: Tanzania)
			"UA": "Ukraine	1974	.ua	ISO 3166-2:UA	Previous ISO country name: Ukrainian SSR
			"UG": "Uganda	1974	.ug	ISO 3166-2:UG	
			"UM": "United States Minor Outlying Islands	1986	.um	ISO 3166-2:UM	Consists of nine minor insular areas of the United States: Baker Island, Howland Island, Jarvis Island, Johnston Atoll, Kingman Reef, Midway Islands, Navassa Island, Palmyra Atoll, and Wake Island
			"US": "United States	1974	.us	ISO 3166-2:US	
			"UY": "Uruguay	1974	.uy	ISO 3166-2:UY	
			"UZ": "Uzbekistan	1992	.uz	ISO 3166-2:UZ	
			"VA": "Holy See (Vatican City State)	1974	.va	ISO 3166-2:VA	Covers Vatican City, territory of the Holy See
			"VC": "Saint Vincent and the Grenadines	1974	.vc	ISO 3166-2:VC	
			"VE": "Venezuela, Bolivarian Republic of	1974	.ve	ISO 3166-2:VE	ISO country name follows UN designation (common name and previous ISO country name: Venezuela)
			"VG": "Virgin Islands, British	1974	.vg	ISO 3166-2:VG	
			"VI": "Virgin Islands, U.S.	1974	.vi	ISO 3166-2:VI	
			"VN": "Viet Nam	1974	.vn	ISO 3166-2:VN	ISO country name follows UN spelling (common spelling: Vietnam)
			"VU": "Vanuatu	1980	.vu	ISO 3166-2:VU	Name changed from New Hebrides (NH)
			"WF": "Wallis and Futuna	1974	.wf	ISO 3166-2:WF	
			"WS": "Samoa	1974	.ws	ISO 3166-2:WS	Code taken from former name: Western Samoa
			"YE": "Yemen	1974	.ye	ISO 3166-2:YE	Previous ISO country name: Yemen, Republic of
			"YT": "Mayotte	1993	.yt	ISO 3166-2:YT	
			"ZA": "South Africa	1974	.za	ISO 3166-2:ZA	Code taken from name in Dutch: Zuid-Afrika
			"ZM": "Zambia	1974	.zm	ISO 3166-2:ZM	
			"ZW": "Zimbabwe	1980	.zw	ISO 3166-2:ZW	Name changed from Southern Rhodesia (RH)
			*/
		}
	}
	return iso3166_1_alpha2
}

