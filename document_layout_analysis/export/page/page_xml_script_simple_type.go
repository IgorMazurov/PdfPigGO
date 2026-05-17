package page

// PageXmlScriptSimpleType represents ISO 15924 script codes (2016-07-14) for PAGE XML documents.
type PageXmlScriptSimpleType byte

const (
	// AdlmAdlam represents the Adlam script.
	AdlmAdlam PageXmlScriptSimpleType = iota

	// AfakAfaka represents the Afaka script.
	AfakAfaka

	// AghbCaucasianAlbanian represents the Caucasian Albanian script.
	AghbCaucasianAlbanian

	// AhomAhomTaiAhom represents the Ahom, Tai Ahom script.
	AhomAhomTaiAhom

	// ArabArabic represents the Arabic script.
	ArabArabic

	// AranArabicNastaliqVariant represents the Arabic (Nastaliq variant) script.
	AranArabicNastaliqVariant

	// ArmiImperialAramaic represents the Imperial Aramaic script.
	ArmiImperialAramaic

	// ArmnArmenian represents the Armenian script.
	ArmnArmenian

	// AvstAvestan represents the Avestan script.
	AvstAvestan

	// BaliBalinese represents the Balinese script.
	BaliBalinese

	// BamuBamum represents the Bamum script.
	BamuBamum

	// BassBassaVah represents the Bassa Vah script.
	BassBassaVah

	// BatkBatak represents the Batak script.
	BatkBatak

	// BengBengali represents the Bengali script.
	BengBengali

	// BhksBhaiksuki represents the Bhaiksuki script.
	BhksBhaiksuki

	// BlisBlissymbols represents the Blissymbols script.
	BlisBlissymbols

	// BopoBopomofo represents the Bopomofo script.
	BopoBopomofo

	// BrahBrahmi represents the Brahmi script.
	BrahBrahmi

	// BraiBraille represents the Braille script.
	BraiBraille

	// BugiBuginese represents the Buginese script.
	BugiBuginese

	// BuhdBuhid represents the Buhid script.
	BuhdBuhid

	// CakmChakma represents the Chakma script.
	CakmChakma

	// CansUnifiedCanadianAboriginalSyllabics represents the Unified Canadian Aboriginal Syllabics script.
	CansUnifiedCanadianAboriginalSyllabics

	// CariCarian represents the Carian script.
	CariCarian

	// ChamCham represents the Cham script.
	ChamCham

	// CherCherokee represents the Cherokee script.
	CherCherokee

	// CirtCirth represents the Cirth script.
	CirtCirth

	// CoptCoptic represents the Coptic script.
	CoptCoptic

	// CprtCypriot represents the Cypriot script.
	CprtCypriot

	// CyrlCyrillic represents the Cyrillic script.
	CyrlCyrillic

	// CyrsCyrillicOldChurchSlavonicVariant represents the Cyrillic (Old Church Slavonic variant) script.
	CyrsCyrillicOldChurchSlavonicVariant

	// DevaDevanagariNagari represents the Devanagari (Nagari) script.
	DevaDevanagariNagari

	// DsrtDeseretMormon represents the Deseret (Mormon) script.
	DsrtDeseretMormon

	// DuplDuployanShorthandDuployanStenography represents the Duployan shorthand, Duployan stenography script.
	DuplDuployanShorthandDuployanStenography

	// EgydEgyptianDemotic represents the Egyptian demotic script.
	EgydEgyptianDemotic

	// EgyhEgyptianHieratic represents the Egyptian hieratic script.
	EgyhEgyptianHieratic

	// EgypEgyptianHieroglyphs represents the Egyptian hieroglyphs script.
	EgypEgyptianHieroglyphs

	// ElbaElbasan represents the Elbasan script.
	ElbaElbasan

	// EthiEthiopic represents the Ethiopic script.
	EthiEthiopic

	// GeokKhutsuriAsomtavruliAndNuskhuri represents the Khutsuri (Asomtavruli and Nuskhuri) script.
	GeokKhutsuriAsomtavruliAndNuskhuri

	// GeorGeorgianMkhedruli represents the Georgian (Mkhedruli) script.
	GeorGeorgianMkhedruli

	// GlagGlagolitic represents the Glagolitic script.
	GlagGlagolitic

	// GothGothic represents the Gothic script.
	GothGothic

	// GranGrantha represents the Grantha script.
	GranGrantha

	// GrekGreek represents the Greek script.
	GrekGreek

	// GujrGujarati represents the Gujarati script.
	GujrGujarati

	// GuruGurmukhi represents the Gurmukhi script.
	GuruGurmukhi

	// HanbHanwithBopomofo represents the Han with Bopomofo script.
	HanbHanwithBopomofo

	// HangHangul represents the Hangul script.
	HangHangul

	// HaniHanHanziKanjiHanja represents the Han (Hanzi, Kanji, Hanja) script.
	HaniHanHanziKanjiHanja

	// HanoHanunoo represents the Hanunoo (Hanunóo) script.
	HanoHanunoo

	// HansHanSimplifiedVariant represents the Han (Simplified variant) script.
	HansHanSimplifiedVariant

	// HantHanTraditionalVariant represents the Han (Traditional variant) script.
	HantHanTraditionalVariant

	// HatrHatran represents the Hatran script.
	HatrHatran

	// HebrHebrew represents the Hebrew script.
	HebrHebrew

	// HiraHiragana represents the Hiragana script.
	HiraHiragana

	// HluwAnatolianHieroglyphs represents the Anatolian Hieroglyphs script.
	HluwAnatolianHieroglyphs

	// HmngPahawhHmong represents the Pahawh Hmong script.
	HmngPahawhHmong

	// HrktJapaneseSyllabaries represents the Japanese syllabaries script.
	HrktJapaneseSyllabaries

	// HungOldHungarianHungarianRunic represents the Old Hungarian (Hungarian Runic) script.
	HungOldHungarianHungarianRunic

	// IndsIndusHarappan represents the Indus (Harappan) script.
	IndsIndusHarappan

	// ItalOldItalicEtruscanOscanEtc represents the Old Italic (Etruscan, Oscan etc.) script.
	ItalOldItalicEtruscanOscanEtc

	// JamoJamo represents the Jamo script.
	JamoJamo

	// JavaJavanese represents the Javanese script.
	JavaJavanese

	// JpanJapanese represents the Japanese script.
	JpanJapanese

	// JurcJurchen represents the Jurchen script.
	JurcJurchen

	// KaliKayahLi represents the Kayah Li script.
	KaliKayahLi

	// KanaKatakana represents the Katakana script.
	KanaKatakana

	// KharKharoshthi represents the Kharoshthi script.
	KharKharoshthi

	// KhmrKhmer represents the Khmer script.
	KhmrKhmer

	// KhojKhojki represents the Khojki script.
	KhojKhojki

	// KitlKhitanLargescript represents the Khitan large script.
	KitlKhitanLargescript

	// KitsKhitansmallscript represents the Khitan small script.
	KitsKhitansmallscript

	// KndaKannada represents the Kannada script.
	KndaKannada

	// KoreKoreanaliasforHangulHan represents the Korean (alias for Hangul + Han) script.
	KoreKoreanaliasforHangulHan

	// KpelKpelle represents the Kpelle script.
	KpelKpelle

	// KthiKaithi represents the Kaithi script.
	KthiKaithi

	// LanaTaiThamLanna represents the Tai Tham (Lanna) script.
	LanaTaiThamLanna

	// LaooLao represents the Lao script.
	LaooLao

	// LatfLatinFrakturvariant represents the Latin (Fraktur variant) script.
	LatfLatinFrakturvariant

	// LatgLatinGaelicvariant represents the Latin (Gaelic variant) script.
	LatgLatinGaelicvariant

	// LatnLatin represents the Latin script.
	LatnLatin

	// LekeLeke represents the Leke script.
	LekeLeke

	// LepcLepchaRong represents the Lepcha (Róng) script.
	LepcLepchaRong

	// LimbLimbu represents the Limbu script.
	LimbLimbu

	// LinaLinearA represents the Linear A script.
	LinaLinearA

	// LinbLinearB represents the Linear B script.
	LinbLinearB

	// LisuLisuFraser represents the Lisu (Fraser) script.
	LisuLisuFraser

	// LomaLoma represents the Loma script.
	LomaLoma

	// LyciLycian represents the Lycian script.
	LyciLycian

	// LydiLydian represents the Lydian script.
	LydiLydian

	// MahjMahajani represents the Mahajani script.
	MahjMahajani

	// MandMandaicMandaean represents the Mandaic, Mandaean script.
	MandMandaicMandaean

	// ManiManichaean represents the Manichaean script.
	ManiManichaean

	// MarcMarchen represents the Marchen script.
	MarcMarchen

	// MayaMayanhieroglyphs represents the Mayan hieroglyphs script.
	MayaMayanhieroglyphs

	// MendMendeKikakui represents the Mende Kikakui script.
	MendMendeKikakui

	// MercMeroiticCursive represents the Meroitic Cursive script.
	MercMeroiticCursive

	// MeroMeroiticHieroglyphs represents the Meroitic Hieroglyphs script.
	MeroMeroiticHieroglyphs

	// MlymMalayalam represents the Malayalam script.
	MlymMalayalam

	// ModiModiMoḍī represents the Modi, Moḍī script.
	ModiModiMoḍī

	// MongMongolian represents the Mongolian script.
	MongMongolian

	// MoonMoonMooncodeMoonscriptMoontype represents the Moon (Moon code, Moon script, Moon type) script.
	MoonMoonMooncodeMoonscriptMoontype

	// MrooMroMru represents the Mro, Mru script.
	MrooMroMru

	// MteiMeiteiMayekMeitheiMeetei represents the Meitei Mayek (Meithei, Meetei) script.
	MteiMeiteiMayekMeitheiMeetei

	// MultMultani represents the Multani script.
	MultMultani

	// MymrMyanmarBurmese represents the Myanmar (Burmese) script.
	MymrMyanmarBurmese

	// NarbOldNorthArabianAncientNorthArabian represents the Old North Arabian (Ancient North Arabian) script.
	NarbOldNorthArabianAncientNorthArabian

	// NbatNabataean represents the Nabataean script.
	NbatNabataean

	// NewaNewaNewarNewari represents the Newa, Newar, Newari script.
	NewaNewaNewarNewari

	// NkgbNakhiGeba represents the Nakhi Geba script.
	NkgbNakhiGeba

	// NkooNKo represents the N'Ko script.
	NkooNKo

	// NshuNuShu represents the Nüshu script.
	NshuNuShu

	// OgamOgham represents the Ogham script.
	OgamOgham

	// OlckOlChikiOlCemetOlSantali represents the Ol Chiki (Ol Cemet', Ol, Santali) script.
	OlckOlChikiOlCemetOlSantali

	// OrkhOldTurkicOrkhonRunic represents the Old Turkic, Orkhon Runic script.
	OrkhOldTurkicOrkhonRunic

	// OryaOriya represents the Oriya script.
	OryaOriya

	// OsgeOsage represents the Osage script.
	OsgeOsage

	// OsmaOsmanya represents the Osmanya script.
	OsmaOsmanya

	// PalmPalmyrene represents the Palmyrene script.
	PalmPalmyrene

	// PaucPauCinHau represents the Pau Cin Hau script.
	PaucPauCinHau

	// PermOldPermic represents the Old Permic script.
	PermOldPermic

	// PhagPhagspa represents the Phags-pa script.
	PhagPhagspa

	// PhliInscriptionalPahlavi represents the Inscriptional Pahlavi script.
	PhliInscriptionalPahlavi

	// PhlpPsalterPahlavi represents the Psalter Pahlavi script.
	PhlpPsalterPahlavi

	// PhlvBookPahlavi represents the Book Pahlavi script.
	PhlvBookPahlavi

	// PhnxPhoenician represents the Phoenician script.
	PhnxPhoenician

	// PiqdKlingonKLIPiqaD represents the Klingon (KLI pIqaD) script.
	PiqdKlingonKLIPiqaD

	// PlrdMiaoPollard represents the Miao (Pollard) script.
	PlrdMiaoPollard

	// PrtiInscriptionalParthian represents the Inscriptional Parthian script.
	PrtiInscriptionalParthian

	// RjngRejangRedjangKaganga represents the Rejang (Redjang, Kaganga) script.
	RjngRejangRedjangKaganga

	// RoroRongorongo represents the Rongorongo script.
	RoroRongorongo

	// RunrRunic represents the Runic script.
	RunrRunic

	// SamrSamaritan represents the Samaritan script.
	SamrSamaritan

	// SaraSarati represents the Sarati script.
	SaraSarati

	// SarbOldSouthArabian represents the Old South Arabian script.
	SarbOldSouthArabian

	// SaurSaurashtra represents the Saurashtra script.
	SaurSaurashtra

	// SgnwSignWriting represents the SignWriting script.
	SgnwSignWriting

	// ShawShavianShaw represents the Shavian (Shaw) script.
	ShawShavianShaw

	// ShrdSharadaŚāradā represents the Sharada, Śāradā script.
	ShrdSharadaŚāradā

	// SiddSiddham represents the Siddham script.
	SiddSiddham

	// SindKhudawadiSindhi represents the Khudawadi, Sindhi script.
	SindKhudawadiSindhi

	// SinhSinhala represents the Sinhala script.
	SinhSinhala

	// SoraSoraSompeng represents the Sora Sompeng script.
	SoraSoraSompeng

	// SundSundanese represents the Sundanese script.
	SundSundanese

	// SyloSylotiNagri represents the Syloti Nagri script.
	SyloSylotiNagri

	// SyrcSyriac represents the Syriac script.
	SyrcSyriac

	// SyreSyriacEstrangeloVariant represents the Syriac (Estrangelo variant) script.
	SyreSyriacEstrangeloVariant

	// SyrjSyriacWesternVariant represents the Syriac (Western variant) script.
	SyrjSyriacWesternVariant

	// SyrnSyriacEasternVariant represents the Syriac (Eastern variant) script.
	SyrnSyriacEasternVariant

	// TagbTagbanwa represents the Tagbanwa script.
	TagbTagbanwa

	// TakrTakri represents the Takri script.
	TakrTakri

	// TaleTaiLe represents the Tai Le script.
	TaleTaiLe

	// TaluNewTaiLue represents the New Tai Lue script.
	TaluNewTaiLue

	// TamlTamil represents the Tamil script.
	TamlTamil

	// TangTangut represents the Tangut script.
	TangTangut

	// TavtTaiViet represents the Tai Viet script.
	TavtTaiViet

	// TeluTelugu represents the Telugu script.
	TeluTelugu

	// TengTengwar represents the Tengwar script.
	TengTengwar

	// TfngTifinaghBerber represents the Tifinagh (Berber) script.
	TfngTifinaghBerber

	// TglgTagalogBaybayinAlibata represents the Tagalog (Baybayin, Alibata) script.
	TglgTagalogBaybayinAlibata

	// ThaaThaana represents the Thaana script.
	ThaaThaana

	// ThaiThai represents the Thai script.
	ThaiThai

	// TibtTibetan represents the Tibetan script.
	TibtTibetan

	// TirhTirhuta represents the Tirhuta script.
	TirhTirhuta

	// UgarUgaritic represents the Ugaritic script.
	UgarUgaritic

	// VaiiVai represents the Vai script.
	VaiiVai

	// VispVisibleSpeech represents the Visible Speech script.
	VispVisibleSpeech

	// WaraWarangCitiVarangKshiti represents the Warang Citi (Varang Kshiti) script.
	WaraWarangCitiVarangKshiti

	// WoleWoleai represents the Woleai script.
	WoleWoleai

	// XpeoOldPersian represents the Old Persian script.
	XpeoOldPersian

	// XsuxCuneiformSumeroAkkadian represents the Cuneiform, Sumero-Akkadian script.
	XsuxCuneiformSumeroAkkadian

	// YiiiYi represents the Yi script.
	YiiiYi

	// ZinhCodeForInheritedScript represents the Code for inherited script.
	ZinhCodeForInheritedScript

	// ZmthMathematicalNotation represents the Mathematical notation.
	ZmthMathematicalNotation

	// ZsyeSymbolsEmojiVariant represents the Symbols (Emoji variant).
	ZsyeSymbolsEmojiVariant

	// ZsymSymbols represents the Symbols.
	ZsymSymbols

	// ZxxxCodeForUnwrittenDocuments represents the Code for unwritten documents.
	ZxxxCodeForUnwrittenDocuments

	// ZyyyCodeForUndeterminedScript represents the Code for undetermined script.
	ZyyyCodeForUndeterminedScript

	// ZzzzCodeForUncodedScript represents the Code for uncoded script.
	ZzzzCodeForUncodedScript

	// ScriptOther represents an unrecognized or unspecified script.
	ScriptOther
)
