package page

// PageXmlLanguageSimpleType represents ISO 639.x language codes (2016-07-14) for PAGE XML documents.
type PageXmlLanguageSimpleType byte

const (
	// Abkhaz represents the Abkhaz language.
	Abkhaz PageXmlLanguageSimpleType = iota

	// Afar represents the Afar language.
	Afar

	// Afrikaans represents the Afrikaans language.
	Afrikaans

	// Akan represents the Akan language.
	Akan

	// Albanian represents the Albanian language.
	Albanian

	// Amharic represents the Amharic language.
	Amharic

	// Arabic represents the Arabic language.
	Arabic

	// Aragonese represents the Aragonese language.
	Aragonese

	// Armenian represents the Armenian language.
	Armenian

	// Assamese represents the Assamese language.
	Assamese

	// Avaric represents the Avaric language.
	Avaric

	// Avestan represents the Avestan language.
	Avestan

	// Aymara represents the Aymara language.
	Aymara

	// Azerbaijani represents the Azerbaijani language.
	Azerbaijani

	// Bambara represents the Bambara language.
	Bambara

	// Bashkir represents the Bashkir language.
	Bashkir

	// Basque represents the Basque language.
	Basque

	// Belarusian represents the Belarusian language.
	Belarusian

	// Bengali represents the Bengali language.
	Bengali

	// Bihari represents the Bihari languages.
	Bihari

	// Bislama represents the Bislama language.
	Bislama

	// Bosnian represents the Bosnian language.
	Bosnian

	// Breton represents the Breton language.
	Breton

	// Bulgarian represents the Bulgarian language.
	Bulgarian

	// Burmese represents the Burmese language.
	Burmese

	// Cambodian represents the Cambodian language.
	Cambodian

	// Cantonese represents the Cantonese language.
	Cantonese

	// Catalan represents the Catalan language.
	Catalan

	// Chamorro represents the Chamorro language.
	Chamorro

	// Chechen represents the Chechen language.
	Chechen

	// Chichewa represents the Chichewa language.
	Chichewa

	// Chinese represents the Chinese language.
	Chinese

	// Chuvash represents the Chuvash language.
	Chuvash

	// Cornish represents the Cornish language.
	Cornish

	// Corsican represents the Corsican language.
	Corsican

	// Cree represents the Cree languages.
	Cree

	// Croatian represents the Croatian language.
	Croatian

	// Czech represents the Czech language.
	Czech

	// Danish represents the Danish language.
	Danish

	// Divehi represents the Divehi language.
	Divehi

	// Dutch represents the Dutch language.
	Dutch

	// Dzongkha represents the Dzongkha language.
	Dzongkha

	// English represents the English language.
	English

	// Esperanto represents the Esperanto language.
	Esperanto

	// Estonian represents the Estonian language.
	Estonian

	// Ewe represents the Ewe language.
	Ewe

	// Faroese represents the Faroese language.
	Faroese

	// Fijian represents the Fijian language.
	Fijian

	// Finnish represents the Finnish language.
	Finnish

	// French represents the French language.
	French

	// Fula represents the Fula language.
	Fula

	// Gaelic represents the Gaelic languages.
	Gaelic

	// Galician represents the Galician language.
	Galician

	// Ganda represents the Ganda language.
	Ganda

	// Georgian represents the Georgian language.
	Georgian

	// German represents the German language.
	German

	// Greek represents the Greek language.
	Greek

	// Guarani represents the Guarani language.
	Guarani

	// Gujarati represents the Gujarati language.
	Gujarati

	// Haitian represents the Haitian language.
	Haitian

	// Hausa represents the Hausa language.
	Hausa

	// Hebrew represents the Hebrew language.
	Hebrew

	// Herero represents the Herero language.
	Herero

	// Hindi represents the Hindi language.
	Hindi

	// HiriMotu represents the Hiri Motu language.
	HiriMotu

	// Hungarian represents the Hungarian language.
	Hungarian

	// Icelandic represents the Icelandic language.
	Icelandic

	// Ido represents the Ido language.
	Ido

	// Igbo represents the Igbo language.
	Igbo

	// Indonesian represents the Indonesian language.
	Indonesian

	// Interlingua represents the Interlingua language.
	Interlingua

	// Interlingue represents the Interlingue language.
	Interlingue

	// Inuktitut represents the Inuktitut language.
	Inuktitut

	// Inupiaq represents the Inupiaq language.
	Inupiaq

	// Irish represents the Irish language.
	Irish

	// Italian represents the Italian language.
	Italian

	// Japanese represents the Japanese language.
	Japanese

	// Javanese represents the Javanese language.
	Javanese

	// Kalaallisut represents the Kalaallisut language.
	Kalaallisut

	// Kannada represents the Kannada language.
	Kannada

	// Kanuri represents the Kanuri language.
	Kanuri

	// Kashmiri represents the Kashmiri language.
	Kashmiri

	// Kazakh represents the Kazakh language.
	Kazakh

	// Khmer represents the Khmer language.
	Khmer

	// Kikuyu represents the Kikuyu language.
	Kikuyu

	// Kinyarwanda represents the Kinyarwanda language.
	Kinyarwanda

	// Kirundi represents the Kirundi language.
	Kirundi

	// Komi represents the Komi language.
	Komi

	// Kongo represents the Kongo language.
	Kongo

	// Korean represents the Korean language.
	Korean

	// Kurdish represents the Kurdish language.
	Kurdish

	// Kwanyama represents the Kwanyama language.
	Kwanyama

	// Kyrgyz represents the Kyrgyz language.
	Kyrgyz

	// Lao represents the Lao language.
	Lao

	// Latin represents the Latin language.
	Latin

	// Latvian represents the Latvian language.
	Latvian

	// Limburgish represents the Limburgish language.
	Limburgish

	// Lingala represents the Lingala language.
	Lingala

	// Lithuanian represents the Lithuanian language.
	Lithuanian

	// LubaKatanga represents the Luba-Katanga language.
	LubaKatanga

	// Luxembourgish represents the Luxembourgish language.
	Luxembourgish

	// Macedonian represents the Macedonian language.
	Macedonian

	// Malagasy represents the Malagasy language.
	Malagasy

	// Malay represents the Malay language.
	Malay

	// Malayalam represents the Malayalam language.
	Malayalam

	// Maltese represents the Maltese language.
	Maltese

	// Manx represents the Manx language.
	Manx

	// Maori represents the Māori language.
	Maori

	// Marathi represents the Marathi language.
	Marathi

	// Marshallese represents the Marshallese language.
	Marshallese

	// Mongolian represents the Mongolian language.
	Mongolian

	// Nauru represents the Nauru language.
	Nauru

	// Navajo represents the Navajo language.
	Navajo

	// Ndonga represents the Ndonga language.
	Ndonga

	// Nepali represents the Nepali language.
	Nepali

	// NorthNdebele represents the North Ndebele language.
	NorthNdebele

	// NorthernSami represents the Northern Sami language.
	NorthernSami

	// Norwegian represents the Norwegian language.
	Norwegian

	// NorwegianBokmal represents the Norwegian Bokmål language.
	NorwegianBokmal

	// NorwegianNynorsk represents the Norwegian Nynorsk language.
	NorwegianNynorsk

	// Nuosu represents the Nuosu language.
	Nuosu

	// Occitan represents the Occitan language.
	Occitan

	// Ojibwe represents the Ojibwe language.
	Ojibwe

	// OldChurchSlavonic represents the Old Church Slavonic language.
	OldChurchSlavonic

	// Oriya represents the Oriya language.
	Oriya

	// Oromo represents the Oromo language.
	Oromo

	// Ossetian represents the Ossetian language.
	Ossetian

	// Pali represents the Pāli language.
	Pali

	// Panjabi represents the Panjabi language.
	Panjabi

	// Pashto represents the Pashto language.
	Pashto

	// Persian represents the Persian language.
	Persian

	// Polish represents the Polish language.
	Polish

	// Portuguese represents the Portuguese language.
	Portuguese

	// Punjabi represents the Punjabi language.
	Punjabi

	// Quechua represents the Quechua language.
	Quechua

	// Romanian represents the Romanian language.
	Romanian

	// Romansh represents the Romansh language.
	Romansh

	// Russian represents the Russian language.
	Russian

	// Samoan represents the Samoan language.
	Samoan

	// Sango represents the Sango language.
	Sango

	// Sanskrit represents the Sanskrit language.
	Sanskrit

	// Sardinian represents the Sardinian language.
	Sardinian

	// Serbian represents the Serbian language.
	Serbian

	// Shona represents the Shona language.
	Shona

	// Sindhi represents the Sindhi language.
	Sindhi

	// Sinhala represents the Sinhala language.
	Sinhala

	// Slovak represents the Slovak language.
	Slovak

	// Slovene represents the Slovene language.
	Slovene

	// Somali represents the Somali language.
	Somali

	// SouthNdebele represents the South Ndebele language.
	SouthNdebele

	// SouthernSotho represents the Southern Sotho language.
	SouthernSotho

	// Spanish represents the Spanish language.
	Spanish

	// Sundanese represents the Sundanese language.
	Sundanese

	// Swahili represents the Swahili language.
	Swahili

	// Swati represents the Swati language.
	Swati

	// Swedish represents the Swedish language.
	Swedish

	// Tagalog represents the Tagalog language.
	Tagalog

	// Tahitian represents the Tahitian language.
	Tahitian

	// Tajik represents the Tajik language.
	Tajik

	// Tamil represents the Tamil language.
	Tamil

	// Tatar represents the Tatar language.
	Tatar

	// Telugu represents the Telugu language.
	Telugu

	// Thai represents the Thai language.
	Thai

	// Tibetan represents the Tibetan language.
	Tibetan

	// Tigrinya represents the Tigrinya language.
	Tigrinya

	// Tonga represents the Tonga language.
	Tonga

	// Tsonga represents the Tsonga language.
	Tsonga

	// Tswana represents the Tswana language.
	Tswana

	// Turkish represents the Turkish language.
	Turkish

	// Turkmen represents the Turkmen language.
	Turkmen

	// Twi represents the Twi language.
	Twi

	// Uighur represents the Uighur language.
	Uighur

	// Ukrainian represents the Ukrainian language.
	Ukrainian

	// Urdu represents the Urdu language.
	Urdu

	// Uzbek represents the Uzbek language.
	Uzbek

	// Venda represents the Venda language.
	Venda

	// Vietnamese represents the Vietnamese language.
	Vietnamese

	// Volapuk represents the Volapük language.
	Volapuk

	// Walloon represents the Walloon language.
	Walloon

	// Welsh represents the Welsh language.
	Welsh

	// WesternFrisian represents the Western Frisian language.
	WesternFrisian

	// Wolof represents the Wolof language.
	Wolof

	// Xhosa represents the Xhosa language.
	Xhosa

	// Yiddish represents the Yiddish language.
	Yiddish

	// Yoruba represents the Yoruba language.
	Yoruba

	// Zhuang represents the Zhuang language.
	Zhuang

	// Zulu represents the Zulu language.
	Zulu

	// LanguageOther represents an unrecognized or unspecified language.
	LanguageOther
)
