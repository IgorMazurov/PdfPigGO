package geometry

import (
	"sort"
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
)

const pointTol = 0.001

func pdfPointEqual(a, b core.PdfPoint) bool {
	dx := a.X - b.X
	dy := a.Y - b.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx <= pointTol && dy <= pointTol
}

func pdfPointsContain(hull []core.PdfPoint, target core.PdfPoint) bool {
	for _, p := range hull {
		if pdfPointEqual(p, target) {
			return true
		}
	}
	return false
}

func sortPoints(pts []core.PdfPoint) {
	sort.Slice(pts, func(a, b int) bool {
		if pts[a].X != pts[b].X {
			return pts[a].X < pts[b].X
		}
		return pts[a].Y < pts[b].Y
	})
}

func TestPdfPointOriginIsZero(t *testing.T) {
	origin := core.Origin
	if origin.X != 0 || origin.Y != 0 {
		t.Errorf("expected zero point, got (%g, %g)", origin.X, origin.Y)
	}
}

func TestPdfPointIntsSetValue(t *testing.T) {
	p := core.NewPdfPoint(256, 372)
	if p.X != 256 || p.Y != 372 {
		t.Errorf("got (%g, %g), want (256, 372)", p.X, p.Y)
	}
}

func TestPdfPointDoublesSetValue(t *testing.T) {
	p := core.NewPdfPoint(0.534436, 0.32552)
	if p.X != 0.534436 || p.Y != 0.32552 {
		t.Errorf("got (%g, %g), want (0.534436, 0.32552)", p.X, p.Y)
	}
}

func TestGrahamScanGenerated(t *testing.T) {
	tests := []struct {
		points []core.PdfPoint
		expect []core.PdfPoint
	}{
		{points: []core.PdfPoint{core.NewPdfPoint(374.54011885, 950.71430641), core.NewPdfPoint(731.99394181, 598.6584842), core.NewPdfPoint(156.01864044, 155.99452034), core.NewPdfPoint(58.08361217, 866.17614577), core.NewPdfPoint(601.11501174, 708.0725778), core.NewPdfPoint(20.5844943, 969.90985216), core.NewPdfPoint(832.4426408, 212.33911068), core.NewPdfPoint(181.82496721, 183.40450985), core.NewPdfPoint(304.24224296, 524.75643163), core.NewPdfPoint(431.94501864, 291.2291402), core.NewPdfPoint(611.85289472, 139.49386065), core.NewPdfPoint(292.14464854, 366.36184329), core.NewPdfPoint(456.06998422, 785.17596139), core.NewPdfPoint(199.67378216, 514.23443841), core.NewPdfPoint(592.41456886, 46.45041272), core.NewPdfPoint(607.5448519, 170.52412369), core.NewPdfPoint(65.05159299, 948.88553725), core.NewPdfPoint(965.63203307, 808.39734812), core.NewPdfPoint(304.61376917, 97.67211401), core.NewPdfPoint(684.23302651, 440.15249374)}, expect: []core.PdfPoint{core.NewPdfPoint(374.54011885, 950.71430641), core.NewPdfPoint(156.01864044, 155.99452034), core.NewPdfPoint(20.5844943, 969.90985216), core.NewPdfPoint(832.4426408, 212.33911068), core.NewPdfPoint(592.41456886, 46.45041272), core.NewPdfPoint(965.63203307, 808.39734812), core.NewPdfPoint(304.61376917, 97.67211401)}},
		{points: []core.PdfPoint{core.NewPdfPoint(15.45661653, 928.31856259), core.NewPdfPoint(428.18414832, 966.65481904), core.NewPdfPoint(963.61997709, 853.00945547)}, expect: []core.PdfPoint{core.NewPdfPoint(15.45661653, 928.31856259), core.NewPdfPoint(428.18414832, 966.65481904), core.NewPdfPoint(963.61997709, 853.00945547)}},
		{points: []core.PdfPoint{core.NewPdfPoint(511.34239886, 501.51629469), core.NewPdfPoint(798.29517897, 649.96393078), core.NewPdfPoint(701.96687726, 795.79266944), core.NewPdfPoint(890.00534182, 337.99515685), core.NewPdfPoint(375.58295264, 93.98193984), core.NewPdfPoint(578.280141, 35.9422738), core.NewPdfPoint(465.59801813, 542.64463471), core.NewPdfPoint(286.54125213, 590.83326057), core.NewPdfPoint(30.50024994, 37.34818875), core.NewPdfPoint(822.60056066, 360.19064141)}, expect: []core.PdfPoint{core.NewPdfPoint(798.29517897, 649.96393078), core.NewPdfPoint(701.96687726, 795.79266944), core.NewPdfPoint(890.00534182, 337.99515685), core.NewPdfPoint(578.280141, 35.9422738), core.NewPdfPoint(286.54125213, 590.83326057), core.NewPdfPoint(30.50024994, 37.34818875)}},
		{points: []core.PdfPoint{core.NewPdfPoint(744.3155504888274, 444.2270145755189), core.NewPdfPoint(569.808503219366, 795.0780697464061), core.NewPdfPoint(979.4238467088927, 747.7769460740175), core.NewPdfPoint(263.9656530143659, 534.1164778047929), core.NewPdfPoint(700.0199185779105, 59.67088755550021), core.NewPdfPoint(350.4405052982569, 201.5075034147189), core.NewPdfPoint(951.4434324059339, 276.4851544966993), core.NewPdfPoint(221.2620795357345, 889.4493759697666), core.NewPdfPoint(26.40411497910822, 836.0708485933704), core.NewPdfPoint(967.4534816241033, 692.8854748787957)}, expect: []core.PdfPoint{core.NewPdfPoint(979.4238467088927, 747.7769460740175), core.NewPdfPoint(700.0199185779105, 59.67088755550021), core.NewPdfPoint(350.4405052982569, 201.5075034147189), core.NewPdfPoint(951.4434324059339, 276.4851544966993), core.NewPdfPoint(221.2620795357345, 889.4493759697666), core.NewPdfPoint(26.40411497910822, 836.0708485933704)}},
		{points: []core.PdfPoint{core.NewPdfPoint(120.21183721137064, 840.2513067067979), core.NewPdfPoint(415.52861888639114, 204.0116313873851), core.NewPdfPoint(415.9664683980775, 832.443368995516), core.NewPdfPoint(277.74879682552734, 502.12519702578516), core.NewPdfPoint(395.35090532250103, 384.0997616867551), core.NewPdfPoint(26.98104432607229, 654.2675428223525), core.NewPdfPoint(507.0471750863688, 822.9947002774292)}, expect: []core.PdfPoint{core.NewPdfPoint(120.21183721137064, 840.2513067067979), core.NewPdfPoint(415.52861888639114, 204.0116313873851), core.NewPdfPoint(415.9664683980775, 832.443368995516), core.NewPdfPoint(26.98104432607229, 654.2675428223525), core.NewPdfPoint(507.0471750863688, 822.9947002774292)}},
		{points: []core.PdfPoint{core.NewPdfPoint(101.52924170788602, 157.72180559576165), core.NewPdfPoint(89.13969305519842, 437.7399878683488), core.NewPdfPoint(963.8109736953475, 34.8258382181843), core.NewPdfPoint(323.44164253397423, 760.0469522615875), core.NewPdfPoint(553.8798287152455, 601.37868940649), core.NewPdfPoint(24.169651959285442, 205.03651103426347), core.NewPdfPoint(598.6426557134453, 973.1839362109574), core.NewPdfPoint(404.2148279309453, 642.4272597428419), core.NewPdfPoint(425.99787946258795, 235.338843056625), core.NewPdfPoint(733.837562543066, 304.97834592908157), core.NewPdfPoint(253.06825516770635, 639.1969849161718), core.NewPdfPoint(702.043917830561, 241.6302720665372), core.NewPdfPoint(43.233323888316356, 214.65896998517496), core.NewPdfPoint(192.04610854054195, 609.086536570487), core.NewPdfPoint(93.436372304046, 130.4501167748684)}, expect: []core.PdfPoint{core.NewPdfPoint(89.13969305519842, 437.7399878683488), core.NewPdfPoint(963.8109736953475, 34.8258382181843), core.NewPdfPoint(323.44164253397423, 760.0469522615875), core.NewPdfPoint(24.169651959285442, 205.03651103426347), core.NewPdfPoint(598.6426557134453, 973.1839362109574), core.NewPdfPoint(192.04610854054195, 609.086536570487), core.NewPdfPoint(93.436372304046, 130.4501167748684)}},
		{points: []core.PdfPoint{core.NewPdfPoint(440.61717693569, 109.0846453611043), core.NewPdfPoint(306.3112668458863, 236.31457456363236), core.NewPdfPoint(520.2265439015014, 849.1406402603602), core.NewPdfPoint(991.7951837792907, 716.5936071305995), core.NewPdfPoint(210.02069783227006, 396.9897189605061), core.NewPdfPoint(844.0889978638492, 40.8614368982706), core.NewPdfPoint(105.18279221346272, 598.2338764791806), core.NewPdfPoint(481.07852218200696, 791.2951914667244), core.NewPdfPoint(306.45291569755994, 923.936778153581), core.NewPdfPoint(959.6247632028106, 131.810977239897), core.NewPdfPoint(800.8517955224195, 350.74228914407433), core.NewPdfPoint(663.4680772483941, 352.304508375042), core.NewPdfPoint(222.58588801803648, 629.9145459539977), core.NewPdfPoint(158.01097092542204, 266.66678727838047), core.NewPdfPoint(63.28771038884462, 52.73381563522583), core.NewPdfPoint(931.6304534780719, 325.3140887623456), core.NewPdfPoint(781.7202794605689, 932.8040026798695), core.NewPdfPoint(691.9750284379987, 896.4394191971841), core.NewPdfPoint(354.5191195854491, 319.9839791539775), core.NewPdfPoint(177.99172848582566, 831.9975645935383), core.NewPdfPoint(237.7316548429551, 731.2752736686494), core.NewPdfPoint(989.397359729622, 217.92540678980444), core.NewPdfPoint(509.249200783335, 621.5474147516464), core.NewPdfPoint(744.8852601187607, 917.8174495950592), core.NewPdfPoint(308.86394317440534, 345.03042856046795), core.NewPdfPoint(716.4931552373788, 116.87982985630296), core.NewPdfPoint(8.999928995218953, 752.2282844071755), core.NewPdfPoint(294.01152638529305, 979.8820875548174), core.NewPdfPoint(919.862114089302, 218.09062242740873), core.NewPdfPoint(58.10245428844474, 679.4146645720332), core.NewPdfPoint(298.5501112825898, 76.27703572384658), core.NewPdfPoint(642.7813627195759, 506.9457773647057), core.NewPdfPoint(95.40371527671387, 849.060165476576), core.NewPdfPoint(462.81907479511983, 59.56530610812505), core.NewPdfPoint(883.4162795690498, 4.108824339342565), core.NewPdfPoint(255.30517910829244, 890.9400020559517), core.NewPdfPoint(127.06302493344879, 246.8462701575741), core.NewPdfPoint(995.1831247109899, 371.7538001202371), core.NewPdfPoint(189.37270784653683, 888.4231305357099), core.NewPdfPoint(924.0107481155084, 775.4673044677166), core.NewPdfPoint(865.7452979418968, 373.2043310431542), core.NewPdfPoint(409.76929412279594, 192.26266847186992), core.NewPdfPoint(438.8872219529338, 819.6826378850956)}, expect: []core.PdfPoint{core.NewPdfPoint(991.7951837792907, 716.5936071305995), core.NewPdfPoint(959.6247632028106, 131.810977239897), core.NewPdfPoint(63.28771038884462, 52.73381563522583), core.NewPdfPoint(781.7202794605689, 932.8040026798695), core.NewPdfPoint(989.397359729622, 217.92540678980444), core.NewPdfPoint(8.999928995218953, 752.2282844071755), core.NewPdfPoint(294.01152638529305, 979.8820875548174), core.NewPdfPoint(95.40371527671387, 849.060165476576), core.NewPdfPoint(883.4162795690498, 4.108824339342565), core.NewPdfPoint(995.1831247109899, 371.7538001202371)}},

	}
	for i, tc := range tests {
		result, err := GrahamScan(tc.points)
		if err != nil {
			t.Errorf("case %d: error: %v", i, err)
			continue
		}
		sortPoints(result)
		expected := make([]core.PdfPoint, len(tc.expect))
		copy(expected, tc.expect)
		sortPoints(expected)
		if len(result) != len(expected) {
			t.Errorf("case %d: got %d points, want %d", i, len(result), len(expected))
			continue
		}
		for j := range result {
			if !pdfPointEqual(result[j], expected[j]) {
				t.Errorf("case %d point %d: got (%g,%g), want (%g,%g)", i, j, result[j].X, result[j].Y, expected[j].X, expected[j].Y)
			}
		}
	}
}

func TestMinimumAreaRectangleGenerated(t *testing.T) {
	tests := []struct {
		points []core.PdfPoint
		expect []core.PdfPoint
	}{
		{points: []core.PdfPoint{core.NewPdfPoint(16, 40.948296), core.NewPdfPoint(18, 44.093158), core.NewPdfPoint(21, 48.810451), core.NewPdfPoint(30, 62.96233), core.NewPdfPoint(49, 92.838519), core.NewPdfPoint(55, 102.273105), core.NewPdfPoint(60, 110.13526), core.NewPdfPoint(64, 116.424984), core.NewPdfPoint(65, 117.997415), core.NewPdfPoint(68, 122.714708), core.NewPdfPoint(75, 133.721725), core.NewPdfPoint(84, 147.873604), core.NewPdfPoint(86, 151.018466), core.NewPdfPoint(90, 157.30819), core.NewPdfPoint(97, 168.315207), core.NewPdfPoint(99, 171.460069), core.NewPdfPoint(105, 180.894655), core.NewPdfPoint(106, 182.467086), core.NewPdfPoint(110, 188.75681), core.NewPdfPoint(113, 193.474103), core.NewPdfPoint(119, 202.908689), core.NewPdfPoint(121, 206.053551), core.NewPdfPoint(122, 207.625982), core.NewPdfPoint(123, 209.198413)}, expect: []core.PdfPoint{core.NewPdfPoint(16, 40.948296), core.NewPdfPoint(16, 40.948296), core.NewPdfPoint(123, 209.198413), core.NewPdfPoint(123, 209.198413)}},
		{points: []core.PdfPoint{core.NewPdfPoint(10.5114328889726, 1131.19945806204), core.NewPdfPoint(11.1542565096881, 1195.02763445421), core.NewPdfPoint(15.3153242359356, 1608.19441341462), core.NewPdfPoint(15.795577716642, 1655.88043939693), core.NewPdfPoint(16.3319701583886, 1709.14069661821), core.NewPdfPoint(17.1302715938114, 1788.40680195762), core.NewPdfPoint(17.9770489556469, 1872.48624937427), core.NewPdfPoint(22.4037700355884, 2312.03066688486), core.NewPdfPoint(23.7647419889928, 2447.16627034947), core.NewPdfPoint(26.9327398551415, 2761.72771472434), core.NewPdfPoint(28.7600875228236, 2943.17137283512), core.NewPdfPoint(32.8221934632313, 3346.51189445353), core.NewPdfPoint(34.4943972831614, 3512.55078434895), core.NewPdfPoint(34.5363374306664, 3516.7151663763), core.NewPdfPoint(36.0282475627216, 3664.85207361081)}, expect: []core.PdfPoint{core.NewPdfPoint(10.5114328889726, 1131.19945806204), core.NewPdfPoint(10.5114328889726, 1131.19945806204), core.NewPdfPoint(36.0282475627216, 3664.85207361081), core.NewPdfPoint(36.0282475627216, 3664.85207361081)}},
		{points: []core.PdfPoint{core.NewPdfPoint(446.78, 217.9), core.NewPdfPoint(446.78, 228.82), core.NewPdfPoint(446.78, 247.52), core.NewPdfPoint(446.78, 256.84), core.NewPdfPoint(446.78, 301.4), core.NewPdfPoint(446.78, 321.39), core.NewPdfPoint(446.78, 369.08), core.NewPdfPoint(446.78, 387.05), core.NewPdfPoint(446.78, 393.22), core.NewPdfPoint(446.78, 397.29), core.NewPdfPoint(446.78, 463.16), core.NewPdfPoint(446.78, 471.88), core.NewPdfPoint(446.78, 480.13), core.NewPdfPoint(446.78, 495.82), core.NewPdfPoint(446.78, 498.99)}, expect: []core.PdfPoint{core.NewPdfPoint(446.78, 217.9), core.NewPdfPoint(446.78, 217.9), core.NewPdfPoint(446.78, 498.99), core.NewPdfPoint(446.78, 498.99)}},
		{points: []core.PdfPoint{core.NewPdfPoint(220, 208.821), core.NewPdfPoint(258.92, 208.821), core.NewPdfPoint(268.61, 208.821), core.NewPdfPoint(283.56, 208.821), core.NewPdfPoint(312.49, 208.821), core.NewPdfPoint(344.93, 208.821), core.NewPdfPoint(356, 208.821), core.NewPdfPoint(356.06, 208.821), core.NewPdfPoint(366.71, 208.821), core.NewPdfPoint(371.07, 208.821), core.NewPdfPoint(430.95, 208.821), core.NewPdfPoint(445.84, 208.821), core.NewPdfPoint(464.95, 208.821), core.NewPdfPoint(470.19, 208.821), core.NewPdfPoint(498.19, 208.821)}, expect: []core.PdfPoint{core.NewPdfPoint(220, 208.821), core.NewPdfPoint(220, 208.821), core.NewPdfPoint(498.19, 208.821), core.NewPdfPoint(498.19, 208.821)}},
		{points: []core.PdfPoint{core.NewPdfPoint(433.8664544437276, 532.7739491464265), core.NewPdfPoint(653.8805470338659, 817.2121262831644), core.NewPdfPoint(531.2551432261636, 360.68316491741035), core.NewPdfPoint(418.79076902856593, 111.73491145462933)}, expect: []core.PdfPoint{core.NewPdfPoint(653.880547033866, 817.2121262831644), core.NewPdfPoint(461.3184174727883, 100.31182441133716), core.NewPdfPoint(327.3696296297328, 136.2909748899142), core.NewPdfPoint(519.9317591908105, 853.1912767617414)}},
		{points: []core.PdfPoint{core.NewPdfPoint(556.0000554986003, 83.08490003544556), core.NewPdfPoint(413.249212383334, 618.580307645359), core.NewPdfPoint(526.8827917872168, 111.78528983363456), core.NewPdfPoint(220.19687765118178, 377.83114567978316)}, expect: []core.PdfPoint{core.NewPdfPoint(556.0000554986004, 83.08490003544566), core.NewPdfPoint(315.83633925060445, 19.06273944752293), core.NewPdfPoint(173.0854961353382, 554.5581470574363), core.NewPdfPoint(413.24921238333417, 618.580307645359)}},
		{points: []core.PdfPoint{core.NewPdfPoint(14.004896521066401, 809.4941990011544), core.NewPdfPoint(703.95616092419, 970.7474069029789), core.NewPdfPoint(835.551079058811, 661.2654117428186), core.NewPdfPoint(200.4833132016346, 14.581989889114745), core.NewPdfPoint(73.40355670360321, 345.2226372321663)}, expect: []core.PdfPoint{core.NewPdfPoint(945.970699704068, 188.81492303383516), core.NewPdfPoint(153.25937285270737, 3.544961240987477), core.NewPdfPoint(-32.56092111592515, 798.6109846932103), core.NewPdfPoint(760.1504057354355, 983.8809464860581)}},
		{points: []core.PdfPoint{core.NewPdfPoint(14.004896521066401, 809.4941990011544), core.NewPdfPoint(703.95616092419, 970.7474069029789), core.NewPdfPoint(835.551079058811, 661.2654117428186), core.NewPdfPoint(703.95616092419, 970.7474069029789), core.NewPdfPoint(703.95616092419, 970.7474069029789), core.NewPdfPoint(200.4833132016346, 14.581989889114745), core.NewPdfPoint(200.4833132016346, 14.581989889114745), core.NewPdfPoint(73.40355670360321, 345.2226372321663), core.NewPdfPoint(73.40355670360321, 345.2226372321663), core.NewPdfPoint(73.40355670360321, 345.2226372321663)}, expect: []core.PdfPoint{core.NewPdfPoint(945.970699704068, 188.81492303383516), core.NewPdfPoint(153.25937285270737, 3.544961240987477), core.NewPdfPoint(-32.56092111592515, 798.6109846932103), core.NewPdfPoint(760.1504057354355, 983.8809464860581)}},
		{points: []core.PdfPoint{core.NewPdfPoint(737.1041856902102, 648.0900313433699), core.NewPdfPoint(258.83885597639045, 32.15501719959235), core.NewPdfPoint(354.5500618726748, 908.7838113897652), core.NewPdfPoint(867.6475306924474, 47.361938752654595), core.NewPdfPoint(352.7960490248145, 283.67860449564785), core.NewPdfPoint(955.4087841797756, 833.327418435315), core.NewPdfPoint(578.2403790703082, 67.4511148622331), core.NewPdfPoint(722.9995401934759, 407.36102955779796), core.NewPdfPoint(404.3710508165602, 736.1127320695537), core.NewPdfPoint(56.25949705397548, 45.503737933916824)}, expect: []core.PdfPoint{core.NewPdfPoint(956.3312379371301, 841.5886588013605), core.NewPdfPoint(857.4507470077488, -43.95763180386355), core.NewPdfPoint(56.25949705397548, 45.503737933916824), core.NewPdfPoint(155.13998798335692, 931.0500285391407)}},
		{points: []core.PdfPoint{core.NewPdfPoint(37.79696601832361, 984.5348223984452), core.NewPdfPoint(881.8169100214427, 818.3604232045343), core.NewPdfPoint(732.8834668881201, 907.4453370173243), core.NewPdfPoint(142.89010125285918, 183.62301422208304), core.NewPdfPoint(292.4319539617013, 383.1685740348906), core.NewPdfPoint(658.9302664366852, 781.9569855570855), core.NewPdfPoint(501.7748878713084, 321.0551716869758), core.NewPdfPoint(104.96397346166219, 658.8420562657931), core.NewPdfPoint(931.1420702804029, 235.94015835854032), core.NewPdfPoint(489.5915692144058, 835.989512769871)}, expect: []core.PdfPoint{core.NewPdfPoint(931.1420702804029, 235.9401583585403), core.NewPdfPoint(91.18213738971586, 180.19110018351975), core.NewPdfPoint(37.79696601832373, 984.5348223984452), core.NewPdfPoint(877.7568989090107, 1040.2838805734657)}},
		{points: []core.PdfPoint{core.NewPdfPoint(360.50496116131137, 475.0257075944493), core.NewPdfPoint(463.6183053562707, 550.6502767074434), core.NewPdfPoint(719.7530742800635, 986.4011287537438), core.NewPdfPoint(746.9948030700444, 387.8044519192034), core.NewPdfPoint(868.5846874204865, 248.33194807842352), core.NewPdfPoint(485.3109455640756, 39.94793327837021), core.NewPdfPoint(504.8344865133781, 708.8088010613369), core.NewPdfPoint(352.0102119724019, 820.2239288583693), core.NewPdfPoint(28.306241810454267, 579.9713087166957), core.NewPdfPoint(801.671925406638, 886.8351079919669)}, expect: []core.PdfPoint{core.NewPdfPoint(1131.0092157743165, 443.1030548242194), core.NewPdfPoint(485.483254766563, -36.005383122939776), core.NewPdfPoint(28.306241810454186, 579.9713087166958), core.NewPdfPoint(673.8322028182079, 1059.079746663855)}},
		{points: []core.PdfPoint{core.NewPdfPoint(384.5196462100093, 457.53938224300975), core.NewPdfPoint(388.62152246674617, 197.75832374394065), core.NewPdfPoint(40.7654267847265, 320.51104165848284), core.NewPdfPoint(443.45264559634137, 22.11792681360786), core.NewPdfPoint(345.74164027022573, 484.01784165724456), core.NewPdfPoint(453.33094307272717, 441.7802101118389), core.NewPdfPoint(470.9897811308254, 63.67713677117809), core.NewPdfPoint(277.98105707671505, 321.72593673466497), core.NewPdfPoint(447.8012370058249, 358.25102431521026), core.NewPdfPoint(345.94253780510235, 111.25057954480089)}, expect: []core.PdfPoint{core.NewPdfPoint(647.6991334648283, 297.75247084308063), core.NewPdfPoint(443.4526455963414, 22.117926813607856), core.NewPdfPoint(40.7654267847265, 320.5110416584829), core.NewPdfPoint(245.0119146532134, 596.1455856879556)}},

	}
	for i, tc := range tests {
		rect, err := MinimumAreaRectangle(tc.points)
		if err != nil {
			t.Errorf("case %d: error: %v", i, err)
			continue
		}
		corners := []core.PdfPoint{rect.BottomLeft, rect.BottomRight, rect.TopLeft, rect.TopRight}
		sortPoints(corners)
		expected := make([]core.PdfPoint, len(tc.expect))
		copy(expected, tc.expect)
		sortPoints(expected)
		if len(corners) != len(expected) {
			t.Errorf("case %d: got %d corners, want %d", i, len(corners), len(expected))
			continue
		}
		for j := range corners {
			if !pdfPointEqual(corners[j], expected[j]) {
				t.Errorf("case %d corner %d: got (%g,%g), want (%g,%g)", i, j, corners[j].X, corners[j].Y, expected[j].X, expected[j].Y)
			}
		}
	}
}

func TestIssue458Generated(t *testing.T) {
	tests := []struct {
		points []core.PdfPoint
		expect []core.PdfPoint
	}{
		{points: []core.PdfPoint{core.NewPdfPoint(134.74199999999985, 1611.657), core.NewPdfPoint(277.58043749999985, 1611.657), core.NewPdfPoint(314.74199999999985, 1611.657), core.NewPdfPoint(507.0248828124999, 1611.657), core.NewPdfPoint(545.1419999999998, 1611.657), core.NewPdfPoint(632.9658281249997, 1611.657), core.NewPdfPoint(668.5709999999998, 1611.657), core.NewPdfPoint(831.3395078125002, 1611.657), core.NewPdfPoint(868.1139999999998, 1611.657), core.NewPdfPoint(892.8907187499999, 1611.657), core.NewPdfPoint(1010.0569999999999, 1611.6569999999997), core.NewPdfPoint(1046.2174003906248, 1611.7654812011715), core.NewPdfPoint(1010.0569999999999, 1611.657), core.NewPdfPoint(1046.2174003906248, 1611.7654812011717), core.NewPdfPoint(1085.145, 1611.8822639999996), core.NewPdfPoint(1255.484941406251, 1612.3932838242179), core.NewPdfPoint(1301.144, 1612.5302609999994), core.NewPdfPoint(1359.0006406250002, 1612.7038309218747), core.NewPdfPoint(1301.144, 1612.5302609999997), core.NewPdfPoint(1359.0006406250002, 1612.703830921875), core.NewPdfPoint(1400.915, 1612.8295739999996), core.NewPdfPoint(1505.379062499999, 1613.1429661874993), core.NewPdfPoint(1400.915, 1612.8295739999999), core.NewPdfPoint(1505.379062499999, 1613.1429661874995), core.NewPdfPoint(1543.886, 1613.2584869999996), core.NewPdfPoint(1600.9401015625003, 1613.4296493046872), core.NewPdfPoint(1641.597, 1613.5516199999997), core.NewPdfPoint(1764.5577421874998, 1613.9205022265612), core.NewPdfPoint(1641.597, 1613.55162), core.NewPdfPoint(1764.5577421874998, 1613.9205022265614)}, expect: []core.PdfPoint{core.NewPdfPoint(134.74199999999985, 1611.657), core.NewPdfPoint(1010.0569999999999, 1611.6569999999997), core.NewPdfPoint(1764.5577421874998, 1613.9205022265614)}},

	}
	for i, tc := range tests {
		result, err := GrahamScan(tc.points)
		if err != nil {
			t.Errorf("case %d: error: %v", i, err)
			continue
		}
		for _, exp := range tc.expect {
			if !pdfPointsContain(result, exp) {
				t.Errorf("case %d: expected point (%g,%g) not in result", i, exp.X, exp.Y)
				break
			}
		}
	}
}

func TestIssue458InvGenerated(t *testing.T) {
	tests := []struct {
		points []core.PdfPoint
		expect []core.PdfPoint
	}{
		{points: []core.PdfPoint{core.NewPdfPoint(134.74199999999985, 1611.657), core.NewPdfPoint(277.58043749999985, 1611.657), core.NewPdfPoint(314.74199999999985, 1611.657), core.NewPdfPoint(507.0248828124999, 1611.657), core.NewPdfPoint(545.1419999999998, 1611.657), core.NewPdfPoint(632.9658281249997, 1611.657), core.NewPdfPoint(668.5709999999998, 1611.657), core.NewPdfPoint(831.3395078125002, 1611.657), core.NewPdfPoint(868.1139999999998, 1611.657), core.NewPdfPoint(892.8907187499999, 1611.657), core.NewPdfPoint(1010.0569999999999, 1611.6569999999997), core.NewPdfPoint(1046.2174003906248, 1611.7654812011715), core.NewPdfPoint(1010.0569999999999, 1611.657), core.NewPdfPoint(1046.2174003906248, 1611.7654812011717), core.NewPdfPoint(1085.145, 1611.8822639999996), core.NewPdfPoint(1255.484941406251, 1612.3932838242179), core.NewPdfPoint(1301.144, 1612.5302609999994), core.NewPdfPoint(1359.0006406250002, 1612.7038309218747), core.NewPdfPoint(1301.144, 1612.5302609999997), core.NewPdfPoint(1359.0006406250002, 1612.703830921875), core.NewPdfPoint(1400.915, 1612.8295739999996), core.NewPdfPoint(1505.379062499999, 1613.1429661874993), core.NewPdfPoint(1400.915, 1612.8295739999999), core.NewPdfPoint(1505.379062499999, 1613.1429661874995), core.NewPdfPoint(1543.886, 1613.2584869999996), core.NewPdfPoint(1600.9401015625003, 1613.4296493046872), core.NewPdfPoint(1641.597, 1613.5516199999997), core.NewPdfPoint(1764.5577421874998, 1613.9205022265612), core.NewPdfPoint(1641.597, 1613.55162), core.NewPdfPoint(1764.5577421874998, 1613.9205022265614)}, expect: []core.PdfPoint{core.NewPdfPoint(134.74199999999985, 1611.657), core.NewPdfPoint(1010.0569999999999, 1611.6569999999997), core.NewPdfPoint(1764.5577421874998, 1613.9205022265614)}},

	}
	for i, tc := range tests {
		pointsInv := make([]core.PdfPoint, len(tc.points))
		for j, p := range tc.points {
			pointsInv[j] = core.NewPdfPoint(p.Y, p.X)
		}
		expectedInv := make([]core.PdfPoint, len(tc.expect))
		for j, p := range tc.expect {
			expectedInv[j] = core.NewPdfPoint(p.Y, p.X)
		}
		result, err := GrahamScan(pointsInv)
		if err != nil {
			t.Errorf("case %d: error: %v", i, err)
			continue
		}
		for _, exp := range expectedInv {
			if !pdfPointsContain(result, exp) {
				t.Errorf("case %d: expected point (%g,%g) not in result", i, exp.X, exp.Y)
				break
			}
		}
	}
}
