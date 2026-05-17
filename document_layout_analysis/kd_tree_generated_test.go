package document_layout_analysis

import (
	"math"
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
)

func assertNearKdTree(t *testing.T, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("got %.10f, want %.10f", got, want)
	}
}

var tree1Points = []core.PdfPoint{
	core.NewPdfPoint(51, 75),
	core.NewPdfPoint(25, 40),
	core.NewPdfPoint(10, 30),
	core.NewPdfPoint(1, 10),
	core.NewPdfPoint(35, 90),
	core.NewPdfPoint(50, 50),
	core.NewPdfPoint(70, 70),
	core.NewPdfPoint(55, 1),
	core.NewPdfPoint(60, 80),
}

var tree2Points = []core.PdfPoint{
	core.NewPdfPoint(82.45353838109239, 62.415005093558115),
	core.NewPdfPoint(16.445674398013832, 7.051331644986569),
	core.NewPdfPoint(35.583597244228336, 27.46872033901967),
	core.NewPdfPoint(63.8245554804822, 44.89509346225664),
	core.NewPdfPoint(0.7742372570967326, 6.938583013317256),
	core.NewPdfPoint(65.40181734760978, 22.29324720161768),
	core.NewPdfPoint(97.13960454465617, 50.25370334786503),
	core.NewPdfPoint(19.15551535520754, 45.53805856543256),
	core.NewPdfPoint(80.59123120177551, 25.51233282470693),
	core.NewPdfPoint(5.021252133213605, 56.04912906205743),
	core.NewPdfPoint(24.41810676553939, 79.31739454102991),
	core.NewPdfPoint(98.12630738875605, 12.483472098628233),
	core.NewPdfPoint(82.48631571672695, 29.34364176375367),
	core.NewPdfPoint(35.803655214742406, 5.16089154047169),
	core.NewPdfPoint(1.813800375443253, 47.255096977735214),
	core.NewPdfPoint(80.11998397458856, 25.99984035749021),
	core.NewPdfPoint(90.26397977468233, 88.93339150438577),
	core.NewPdfPoint(8.254965696884042, 8.97868520394496),
	core.NewPdfPoint(14.158071542303697, 98.52077265281578),
	core.NewPdfPoint(70.71899728005727, 9.70650341683359),
	core.NewPdfPoint(58.60508497220939, 69.22169822597527),
	core.NewPdfPoint(15.143510163282482, 66.19892308162831),
	core.NewPdfPoint(16.49098129875618, 61.44894011847811),
	core.NewPdfPoint(59.31519405414486, 64.30374393065894),
	core.NewPdfPoint(58.880198617023936, 69.74930498275906),
	core.NewPdfPoint(58.78888045620321, 4.76405876739463),
	core.NewPdfPoint(11.228694290441743, 12.860145308234605),
	core.NewPdfPoint(80.62371302859762, 81.02375037506395),
	core.NewPdfPoint(90.87481745321885, 8.048659530077462),
	core.NewPdfPoint(86.97060905691454, 53.23733886383501),
	core.NewPdfPoint(79.77968787046065, 62.69132229661655),
	core.NewPdfPoint(97.45369858932968, 33.575350634376086),
	core.NewPdfPoint(55.68259729639997, 78.53266468780453),
	core.NewPdfPoint(84.80993603153726, 40.80973096746926),
	core.NewPdfPoint(52.89633110070122, 15.424096971651757),
	core.NewPdfPoint(31.112041569584935, 55.44209995586027),
	core.NewPdfPoint(49.53804959253269, 92.63065004232082),
	core.NewPdfPoint(92.45607820110665, 72.66951564303699),
	core.NewPdfPoint(24.636942156010523, 22.044633809864933),
	core.NewPdfPoint(99.76982890678676, 52.90800974490838),
	core.NewPdfPoint(82.61099784649278, 17.85771383958713),
	core.NewPdfPoint(74.16588992664359, 81.00550169213861),
	core.NewPdfPoint(13.34483441571791, 47.0726836958395),
	core.NewPdfPoint(47.945049322341525, 9.252548969452246),
	core.NewPdfPoint(82.78544384909999, 85.53384162961818),
	core.NewPdfPoint(86.41791572882465, 84.5461317010868),
	core.NewPdfPoint(67.53711275495073, 90.39020510040851),
	core.NewPdfPoint(62.56803790975437, 46.05914229645705),
	core.NewPdfPoint(71.56130305322326, 8.098723357173876),
	core.NewPdfPoint(37.970724110053844, 36.503310666941424),
	core.NewPdfPoint(18.3788032374675, 19.84795034896214),
	core.NewPdfPoint(18.13308334914272, 1.2881803099952127),
	core.NewPdfPoint(81.60160120456186, 59.00242318333139),
	core.NewPdfPoint(37.57869607184197, 27.967241358959505),
	core.NewPdfPoint(97.7460348272766, 5.634112362720456),
	core.NewPdfPoint(84.9666796464414, 68.61319827025679),
	core.NewPdfPoint(29.725739779656745, 38.277277199189065),
	core.NewPdfPoint(15.448731011261907, 30.373061794712797),
	core.NewPdfPoint(70.09436264721477, 86.01779473250726),
	core.NewPdfPoint(71.14819197800337, 52.17424240997747),
	core.NewPdfPoint(43.017566199578106, 57.09185921458071),
	core.NewPdfPoint(33.433656443618766, 3.2699185921492235),
	core.NewPdfPoint(2.4884293540186175, 29.514037079380717),
	core.NewPdfPoint(72.74747856539602, 44.441442049411116),
	core.NewPdfPoint(15.264247773398054, 0.8939481426395335),
	core.NewPdfPoint(82.27200991589356, 79.74763514660499),
	core.NewPdfPoint(34.43836768761709, 85.52303143204219),
	core.NewPdfPoint(91.90691599722378, 96.20420265233045),
	core.NewPdfPoint(10.747687004702266, 72.62390777764313),
	core.NewPdfPoint(18.204503055467892, 62.00619128580146),
	core.NewPdfPoint(91.00755813085259, 81.62926984298593),
	core.NewPdfPoint(52.97863363490906, 68.26721122337617),
	core.NewPdfPoint(84.97188896260866, 40.483497230772926),
	core.NewPdfPoint(82.04208148460329, 71.99786553392418),
	core.NewPdfPoint(75.03458936225277, 45.61844875850991),
	core.NewPdfPoint(62.777729239709636, 93.54972672515291),
	core.NewPdfPoint(70.36073331446983, 3.9743344551363857),
	core.NewPdfPoint(84.64548329605812, 44.580018552016185),
	core.NewPdfPoint(44.866927314279856, 7.669932380951572),
	core.NewPdfPoint(90.03637785608008, 39.71900281528301),
	core.NewPdfPoint(63.98092485758994, 24.420561596263646),
	core.NewPdfPoint(11.201209707090698, 49.55375920684947),
	core.NewPdfPoint(99.66377946738517, 22.268044551668787),
	core.NewPdfPoint(57.19140901930086, 6.016576988006017),
	core.NewPdfPoint(42.195134616764996, 62.54652068319705),
	core.NewPdfPoint(49.61310425442538, 20.010647167526496),
	core.NewPdfPoint(52.44257355334966, 84.19422038812934),
	core.NewPdfPoint(67.0752294529232, 39.24937886844233),
	core.NewPdfPoint(43.03278016511083, 4.6357991376672185),
	core.NewPdfPoint(49.42598643546282, 18.91930873569896),
	core.NewPdfPoint(93.45285655237589, 26.040139241642603),
	core.NewPdfPoint(99.36796037285181, 90.9895050206579),
	core.NewPdfPoint(69.38884149074713, 6.0366785105100185),
	core.NewPdfPoint(99.69679458529599, 47.806517636673476),
	core.NewPdfPoint(80.69475374791655, 21.854199922464968),
	core.NewPdfPoint(89.04717161055198, 88.20721990563409),
	core.NewPdfPoint(46.898436066324166, 25.173429288152636),
	core.NewPdfPoint(8.420107657979937, 68.0080815210619),
	core.NewPdfPoint(28.000497306010118, 81.20161949340638),
	core.NewPdfPoint(24.103817108644833, 21.163820177769633),
}

// TestBuildTree matches C# [Fact] BuildTree()
func TestBuildTree(t *testing.T) {
	candidates := []core.PdfPoint{
		core.NewPdfPoint(2, 3),
		core.NewPdfPoint(4, 7),
		core.NewPdfPoint(5, 4),
		core.NewPdfPoint(7, 2),
		core.NewPdfPoint(8, 1),
		core.NewPdfPoint(9, 6),
	}

	tree, err := NewKdTree(candidates)
	if err != nil { t.Fatal(err) }

	root := tree.Root()

// root checks
	if math.Abs(root.value.X-7) > 0.001 || math.Abs(root.value.Y-2) > 0.001 { t.Errorf("root.value: got %v", root.value) }
	if root.depth != 0 { t.Errorf("root.depth: got %d", root.depth) }
	if root.index != 3 { t.Errorf("root.index: got %d", root.index) }
	if root.isLeaf { t.Error("root.isLeaf should be false") }
	if !root.isAxisCutX { t.Error("root.isAxisCutX should be true") }

// root -> left side
	left := root.leftChild
	if left == nil { t.Fatal("root.leftChild is nil") }
	if math.Abs(left.value.X-5) > 0.001 || math.Abs(left.value.Y-4) > 0.001 { t.Errorf("left.value: got %v", left.value) }
	if left.depth != 1 { t.Errorf("left.depth: got %d", left.depth) }
	if left.index != 2 { t.Errorf("left.index: got %d", left.index) }
	if left.isLeaf { t.Error("left.isLeaf should be false") }
	if left.isAxisCutX { t.Error("left.isAxisCutX should be false") }

// root -> left -> left
	ll := left.leftChild
	if ll == nil { t.Fatal("left.leftChild is nil") }
	if math.Abs(ll.value.X-2) > 0.001 || math.Abs(ll.value.Y-3) > 0.001 { t.Errorf("ll.value: got %v", ll.value) }
	if ll.depth != 2 { t.Errorf("ll.depth: got %d", ll.depth) }
	if ll.index != 0 { t.Errorf("ll.index: got %d", ll.index) }
	if !ll.isLeaf { t.Error("ll.isLeaf should be true") }
	if !ll.isAxisCutX { t.Error("ll.isAxisCutX should be true") }
	if ll.leftChild != nil || ll.rightChild != nil { t.Error("leaf should have no children") }

// root -> left -> right
	lr := left.rightChild
	if lr == nil { t.Fatal("left.rightChild is nil") }
	if math.Abs(lr.value.X-4) > 0.001 || math.Abs(lr.value.Y-7) > 0.001 { t.Errorf("lr.value: got %v", lr.value) }
	if lr.depth != 2 { t.Errorf("lr.depth: got %d", lr.depth) }
	if lr.index != 1 { t.Errorf("lr.index: got %d", lr.index) }
	if !lr.isLeaf { t.Error("lr.isLeaf should be true") }
	if !lr.isAxisCutX { t.Error("lr.isAxisCutX should be true") }
	if lr.leftChild != nil || lr.rightChild != nil { t.Error("leaf should have no children") }

// root -> right side
	right := root.rightChild
	if right == nil { t.Fatal("root.rightChild is nil") }
	if math.Abs(right.value.X-9) > 0.001 || math.Abs(right.value.Y-6) > 0.001 { t.Errorf("right.value: got %v", right.value) }
	if right.depth != 1 { t.Errorf("right.depth: got %d", right.depth) }
	if right.index != 5 { t.Errorf("right.index: got %d", right.index) }
	if right.isLeaf { t.Error("right.isLeaf should be false") }
	if right.isAxisCutX { t.Error("right.isAxisCutX should be false") }

// root -> right -> left
	rl := right.leftChild
	if rl == nil { t.Fatal("right.leftChild is nil") }
	if math.Abs(rl.value.X-8) > 0.001 || math.Abs(rl.value.Y-1) > 0.001 { t.Errorf("rl.value: got %v", rl.value) }
	if rl.depth != 2 { t.Errorf("rl.depth: got %d", rl.depth) }
	if rl.index != 4 { t.Errorf("rl.index: got %d", rl.index) }
	if !rl.isLeaf { t.Error("rl.isLeaf should be true") }
	if !rl.isAxisCutX { t.Error("rl.isAxisCutX should be true") }
	if rl.leftChild != nil || rl.rightChild != nil { t.Error("leaf should have no children") }

// root -> right -> right (should be nil)
	if right.rightChild != nil { t.Error("right.rightChild should be nil") }
}

// 51 cases from C# DataTree1
func TestNearestTree1Generated(t *testing.T) {
	tree, err := NewKdTree(tree1Points)
	if err != nil { t.Fatal(err) }

	tests := []struct{ query core.PdfPoint; expectDist float64; expectIdx int; expectPt core.PdfPoint }{
		{query: core.NewPdfPoint(51, 49), expectDist: 1.4142135623730951, expectIdx: 5, expectPt: core.NewPdfPoint(50, 50)},
		{query: core.NewPdfPoint(28.189524796700038, 75.60283789175995), expectDist: 15.926733791512522, expectIdx: 4, expectPt: core.NewPdfPoint(35, 90)},
		{query: core.NewPdfPoint(43.26688589899484, 8.035369191312736), expectDist: 13.680730468994646, expectIdx: 7, expectPt: core.NewPdfPoint(55, 1)},
		{query: core.NewPdfPoint(82.22662843535518, 70.12992266643707), expectDist: 12.227318708346896, expectIdx: 6, expectPt: core.NewPdfPoint(70, 70)},
		{query: core.NewPdfPoint(76.29751404813953, 74.63310544916789), expectDist: 7.8182062705983855, expectIdx: 6, expectPt: core.NewPdfPoint(70, 70)},
		{query: core.NewPdfPoint(50.895502833937776, 61.76937358125091), expectDist: 11.80339272500231, expectIdx: 5, expectPt: core.NewPdfPoint(50, 50)},
		{query: core.NewPdfPoint(42.9552821543992, 74.01081889822333), expectDist: 8.105304711572543, expectIdx: 0, expectPt: core.NewPdfPoint(51, 75)},
		{query: core.NewPdfPoint(55.51285821663918, 76.33782834155515), expectDist: 4.706981405843449, expectIdx: 0, expectPt: core.NewPdfPoint(51, 75)},
		{query: core.NewPdfPoint(4.3199936890310315, 98.54112120917016), expectDist: 31.846719434673833, expectIdx: 4, expectPt: core.NewPdfPoint(35, 90)},
		{query: core.NewPdfPoint(21.569550153382444, 37.56786718125442), expectDist: 4.205146394381262, expectIdx: 1, expectPt: core.NewPdfPoint(25, 40)},
		{query: core.NewPdfPoint(95.70493339772732, 65.77875848107642), expectDist: 26.04923186857304, expectIdx: 6, expectPt: core.NewPdfPoint(70, 70)},
		{query: core.NewPdfPoint(87.23320341806003, 56.082576505219414), expectDist: 22.151252262147768, expectIdx: 6, expectPt: core.NewPdfPoint(70, 70)},
		{query: core.NewPdfPoint(95.64105363103229, 87.41037179209023), expectDist: 30.99330052202065, expectIdx: 6, expectPt: core.NewPdfPoint(70, 70)},
		{query: core.NewPdfPoint(31.9581372373153, 85.00443296498887), expectDist: 5.848813475252715, expectIdx: 4, expectPt: core.NewPdfPoint(35, 90)},
		{query: core.NewPdfPoint(36.17227123238111, 79.38086715887763), expectDist: 10.68364180135556, expectIdx: 4, expectPt: core.NewPdfPoint(35, 90)},
		{query: core.NewPdfPoint(97.45057198438961, 0.1038719192212767), expectDist: 42.46002952588474, expectIdx: 7, expectPt: core.NewPdfPoint(55, 1)},
		{query: core.NewPdfPoint(49.936342420193036, 26.674269408896535), expectDist: 23.3258174539759, expectIdx: 5, expectPt: core.NewPdfPoint(50, 50)},
		{query: core.NewPdfPoint(96.58603550736572, 5.478651527669465), expectDist: 41.826506771737215, expectIdx: 7, expectPt: core.NewPdfPoint(55, 1)},
		{query: core.NewPdfPoint(33.82105876943279, 99.91438503275373), expectDist: 9.984234222153567, expectIdx: 4, expectPt: core.NewPdfPoint(35, 90)},
		{query: core.NewPdfPoint(40.95742737577155, 81.48724148079593), expectDist: 10.390283852901893, expectIdx: 4, expectPt: core.NewPdfPoint(35, 90)},
		{query: core.NewPdfPoint(65.32548187684739, 91.29488761267328), expectDist: 12.487403389157823, expectIdx: 8, expectPt: core.NewPdfPoint(60, 80)},
		{query: core.NewPdfPoint(95.72487295050016, 17.00011169070058), expectDist: 43.75521512859093, expectIdx: 7, expectPt: core.NewPdfPoint(55, 1)},
		{query: core.NewPdfPoint(76.84590618203656, 90.69206464298559), expectDist: 19.952563780721515, expectIdx: 8, expectPt: core.NewPdfPoint(60, 80)},
		{query: core.NewPdfPoint(27.727148118788257, 46.59561305641292), expectDist: 7.137187713079637, expectIdx: 1, expectPt: core.NewPdfPoint(25, 40)},
		{query: core.NewPdfPoint(51.92233547390256, 48.53466193663921), expectDist: 2.4171448682605097, expectIdx: 5, expectPt: core.NewPdfPoint(50, 50)},
		{query: core.NewPdfPoint(74.85516050981272, 96.58650922810423), expectDist: 22.26629924676048, expectIdx: 8, expectPt: core.NewPdfPoint(60, 80)},
		{query: core.NewPdfPoint(83.16610264969573, 97.29274249477987), expectDist: 28.908601747006117, expectIdx: 8, expectPt: core.NewPdfPoint(60, 80)},
		{query: core.NewPdfPoint(71.98380721011306, 96.58511553469276), expectDist: 20.46161510116601, expectIdx: 8, expectPt: core.NewPdfPoint(60, 80)},
		{query: core.NewPdfPoint(98.51219383024967, 32.70778133798633), expectDist: 46.943101407439705, expectIdx: 6, expectPt: core.NewPdfPoint(70, 70)},
		{query: core.NewPdfPoint(8.09276417067899, 51.21666984962398), expectDist: 20.28961078738919, expectIdx: 1, expectPt: core.NewPdfPoint(25, 40)},
		{query: core.NewPdfPoint(61.41570675706768, 78.62680876706933), expectDist: 1.9722778161822772, expectIdx: 8, expectPt: core.NewPdfPoint(60, 80)},
		{query: core.NewPdfPoint(54.00356297008202, 36.88821236457238), expectDist: 13.709394277354654, expectIdx: 5, expectPt: core.NewPdfPoint(50, 50)},
		{query: core.NewPdfPoint(21.85770655063177, 16.016459205926083), expectDist: 18.334247128814017, expectIdx: 2, expectPt: core.NewPdfPoint(10, 30)},
		{query: core.NewPdfPoint(27.683193240146863, 55.85996630126125), expectDist: 16.085336708975426, expectIdx: 1, expectPt: core.NewPdfPoint(25, 40)},
		{query: core.NewPdfPoint(65.54133757456142, 92.08898521546233), expectDist: 13.298495616230923, expectIdx: 8, expectPt: core.NewPdfPoint(60, 80)},
		{query: core.NewPdfPoint(24.861338264089227, 20.46892423497547), expectDist: 17.65504970931321, expectIdx: 2, expectPt: core.NewPdfPoint(10, 30)},
		{query: core.NewPdfPoint(54.88661497017837, 4.925311272555266), expectDist: 3.92694852925743, expectIdx: 7, expectPt: core.NewPdfPoint(55, 1)},
		{query: core.NewPdfPoint(27.79457376217088, 54.31806138769485), expectDist: 14.58823239511943, expectIdx: 1, expectPt: core.NewPdfPoint(25, 40)},
		{query: core.NewPdfPoint(0.7501079718312487, 80.59920809402051), expectDist: 35.51661572279581, expectIdx: 4, expectPt: core.NewPdfPoint(35, 90)},
		{query: core.NewPdfPoint(89.56921087453362, 16.96508157736404), expectDist: 38.07773851294037, expectIdx: 7, expectPt: core.NewPdfPoint(55, 1)},
		{query: core.NewPdfPoint(5.3978925878381485, 59.8520804896729), expectDist: 27.898883754845496, expectIdx: 1, expectPt: core.NewPdfPoint(25, 40)},
		{query: core.NewPdfPoint(25.97829935047318, 1.9752285215369203), expectDist: 26.23570840902535, expectIdx: 3, expectPt: core.NewPdfPoint(1, 10)},
		{query: core.NewPdfPoint(62.572023839684796, 39.52157530571068), expectDist: 16.366220318066574, expectIdx: 5, expectPt: core.NewPdfPoint(50, 50)},
		{query: core.NewPdfPoint(81.17046822810447, 92.91856347287415), expectDist: 24.800766262352845, expectIdx: 8, expectPt: core.NewPdfPoint(60, 80)},
		{query: core.NewPdfPoint(72.71721419976534, 53.943166291078306), expectDist: 16.285120870394863, expectIdx: 6, expectPt: core.NewPdfPoint(70, 70)},
		{query: core.NewPdfPoint(77.0888492617333, 55.48826101878681), expectDist: 16.150614604851395, expectIdx: 6, expectPt: core.NewPdfPoint(70, 70)},
		{query: core.NewPdfPoint(58.21541870080057, 11.40938709274073), expectDist: 10.894689397498917, expectIdx: 7, expectPt: core.NewPdfPoint(55, 1)},
		{query: core.NewPdfPoint(92.47344710876337, 10.739546752843566), expectDist: 38.718445335061055, expectIdx: 7, expectPt: core.NewPdfPoint(55, 1)},
		{query: core.NewPdfPoint(98.99910001507719, 94.33068443488959), expectDist: 37.854061958456036, expectIdx: 6, expectPt: core.NewPdfPoint(70, 70)},
		{query: core.NewPdfPoint(0.29975539175036703, 35.105628826864496), expectDist: 10.961851630887265, expectIdx: 2, expectPt: core.NewPdfPoint(10, 30)},
		{query: core.NewPdfPoint(69.86806543596909, 86.99952741106358), expectDist: 12.09843375924331, expectIdx: 8, expectPt: core.NewPdfPoint(60, 80)},
	}
	for _, tc := range tests {
		point, idx, dist := tree.FindNearestNeighbour(tc.query, Euclidean)
		assertNearKdTree(t, dist, tc.expectDist, 0.000001)
		if idx != tc.expectIdx { t.Errorf("idx: got %d, want %d", idx, tc.expectIdx) }
		if math.Abs(point.X-tc.expectPt.X)>0.001||math.Abs(point.Y-tc.expectPt.Y)>0.001 {
			t.Errorf("pt: got %v, want %v", point, tc.expectPt)
		}
	}
}

// 100 cases from C# DataTreeK1 (k=3)
func TestNearestKTree1Generated(t *testing.T) {
	tree, err := NewKdTree(tree1Points)
	if err != nil { t.Fatal(err) }

	tests := []struct{ query core.PdfPoint; expected [3]NearestResult }{
		{query: core.NewPdfPoint(57.28490719962775, 71.58006710449683), expected: [3]NearestResult{{Distance: 7.155137980338146, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 8.846863787773017, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 12.812891827249272, Index: 6, Point: core.NewPdfPoint(70, 70)}}},
		{query: core.NewPdfPoint(57.6366179124611, 99.23965351255829), expected: [3]NearestResult{{Distance: 19.38426790402455, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 24.449696675968916, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 25.13176276596767, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(73.90499795029336, 90.88238489014321), expected: [3]NearestResult{{Distance: 17.657159139988515, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 21.24436413950479, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 27.87273005827726, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(23.565146028832896, 75.26741657423868), expected: [3]NearestResult{{Distance: 18.64952813716565, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 27.43615723900563, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 35.29659300470031, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(26.939537622638365, 12.273123896672733), expected: [3]NearestResult{{Distance: 24.51917762184319, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 26.038945914262655, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 27.79463014035068, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(83.86699657545552, 96.41621467187889), expected: [3]NearestResult{{Distance: 28.967665243958084, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 29.83471118704662, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 39.228735829274505, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(66.43635368261755, 7.7006684896940625), expected: [3]NearestResult{{Distance: 13.254778148377245, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 45.380471224953766, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 52.537779002576954, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(37.3870248673188, 57.9719767606922), expected: [3]NearestResult{{Distance: 14.921111056843282, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 21.800611629066815, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 21.827284158832835, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(28.21563125325024, 51.82289295570881), expected: [3]NearestResult{{Distance: 12.252390876846393, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 21.860504578402132, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 28.42618298172852, Index: 2, Point: core.NewPdfPoint(10, 30)}}},
		{query: core.NewPdfPoint(40.64586336533863, 54.586928359658906), expected: [3]NearestResult{{Distance: 10.418242843999996, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 21.39092142514358, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 22.888897728875776, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(22.971155983334967, 41.71022847164334), expected: [3]NearestResult{{Distance: 2.6535051289147757, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 17.475134860769824, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 28.271517838092137, Index: 5, Point: core.NewPdfPoint(50, 50)}}},
		{query: core.NewPdfPoint(46.16350096686003, 0.026010364849260448), expected: [3]NearestResult{{Distance: 8.890015240260544, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 45.230671236732974, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 46.25172741450487, Index: 3, Point: core.NewPdfPoint(1, 10)}}},
		{query: core.NewPdfPoint(30.810311156595372, 57.2429846701653), expected: [3]NearestResult{{Distance: 18.195610351730775, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 20.51109418921715, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 26.88745300353832, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(8.49831703999272, 7.098709173582418), expected: [3]NearestResult{{Distance: 8.040040229482686, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 22.95047217877084, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 36.80761441002533, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(45.59952054686074, 21.5904374552803), expected: [3]NearestResult{{Distance: 22.634821151241802, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 27.62702010439208, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 28.74834714205047, Index: 5, Point: core.NewPdfPoint(50, 50)}}},
		{query: core.NewPdfPoint(75.87075134690805, 58.59787314401566), expected: [3]NearestResult{{Distance: 12.824750220459736, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 26.644545075400142, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 27.262046839042196, Index: 5, Point: core.NewPdfPoint(50, 50)}}},
		{query: core.NewPdfPoint(8.963229075823397, 94.18880408311563), expected: [3]NearestResult{{Distance: 26.37156650267053, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 46.20930979653228, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 52.972390428183594, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(57.181014731701026, 9.603539673394279), expected: [3]NearestResult{{Distance: 8.875681391958942, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 41.02975724393135, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 44.26694601560936, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(54.645113285069115, 27.785970730398834), expected: [3]NearestResult{{Distance: 22.694496553610147, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 26.788321570233126, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 32.062676941940694, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(11.149446927742158, 88.54384340073094), expected: [3]NearestResult{{Distance: 23.89496335829337, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 42.089218028235706, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 49.59207391833591, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(99.50678603018885, 74.5426364401791), expected: [3]NearestResult{{Distance: 29.854412867430348, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 39.881937759582165, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 48.50894219011971, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(24.209348883557468, 41.95277654538159), expected: [3]NearestResult{{Distance: 2.106766580360595, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 18.568103372140087, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 27.016948205499066, Index: 5, Point: core.NewPdfPoint(50, 50)}}},
		{query: core.NewPdfPoint(71.52714336370258, 2.423670062161676), expected: [3]NearestResult{{Distance: 16.58834844733717, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 52.21996813246304, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 59.805983322605925, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(64.29355652488779, 17.6151418970531), expected: [3]NearestResult{{Distance: 19.037676673914117, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 35.39893773092871, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 45.222399943652036, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(2.3569761133599876, 34.0132529249903), expected: [3]NearestResult{{Distance: 8.632613345429817, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 23.42109457884254, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 24.051563363153438, Index: 3, Point: core.NewPdfPoint(1, 10)}}},
		{query: core.NewPdfPoint(47.10705238689328, 14.985277112050543), expected: [3]NearestResult{{Distance: 16.058847963789056, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 33.383500810995635, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 35.134028587852995, Index: 5, Point: core.NewPdfPoint(50, 50)}}},
		{query: core.NewPdfPoint(73.44870263360212, 44.499915917034315), expected: [3]NearestResult{{Distance: 24.08511117098677, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 25.73223344549272, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 37.87082490519413, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(37.15926087979008, 21.563350867442598), expected: [3]NearestResult{{Distance: 22.08523616309826, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 27.223948487553027, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 28.439453231776056, Index: 2, Point: core.NewPdfPoint(10, 30)}}},
		{query: core.NewPdfPoint(19.784956200369862, 49.74041414042939), expected: [3]NearestResult{{Distance: 11.048635637902876, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 22.032460558884953, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 30.216158866276444, Index: 5, Point: core.NewPdfPoint(50, 50)}}},
		{query: core.NewPdfPoint(91.46918827031537, 10.458352581834374), expected: [3]NearestResult{{Distance: 37.675749848649346, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 57.29952404986787, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 63.29402675020287, Index: 6, Point: core.NewPdfPoint(70, 70)}}},
		{query: core.NewPdfPoint(19.900141743444067, 33.68929824559639), expected: [3]NearestResult{{Distance: 8.113785236866608, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 10.565213111208138, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 30.305085535122533, Index: 3, Point: core.NewPdfPoint(1, 10)}}},
		{query: core.NewPdfPoint(99.56220935020619, 63.18918949351682), expected: [3]NearestResult{{Distance: 30.3366339830351, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 42.985715750170165, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 49.97782930253481, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(0.9575689831946121, 56.0195244568354), expected: [3]NearestResult{{Distance: 27.545983584790353, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 28.890546083814222, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 46.01954401799804, Index: 3, Point: core.NewPdfPoint(1, 10)}}},
		{query: core.NewPdfPoint(13.763956851598868, 20.74365587257857), expected: [3]NearestResult{{Distance: 9.992360971559586, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 16.68366674379076, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 22.294740518481255, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(89.66047561920547, 70.74342301744882), expected: [3]NearestResult{{Distance: 19.67452615328373, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 31.071337779236003, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 38.894097530493816, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(51.22041417582528, 90.11643810630143), expected: [3]NearestResult{{Distance: 13.39490377728326, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 15.118044960594169, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 16.220832095423244, Index: 4, Point: core.NewPdfPoint(35, 90)}}},
		{query: core.NewPdfPoint(13.94537026257473, 3.2534322160761575), expected: [3]NearestResult{{Distance: 14.597903551477287, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 27.035991469314418, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 38.37336423262986, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(35.907000931452906, 67.29871186609165), expected: [3]NearestResult{{Distance: 16.94427513364443, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 22.31273302336873, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 22.719399939883633, Index: 4, Point: core.NewPdfPoint(35, 90)}}},
		{query: core.NewPdfPoint(53.72385478058551, 34.1117586202502), expected: [3]NearestResult{{Distance: 16.318802301887334, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 29.321173578190265, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 33.13634116112924, Index: 7, Point: core.NewPdfPoint(55, 1)}}},
		{query: core.NewPdfPoint(56.01158219636303, 50.713431740553794), expected: [3]NearestResult{{Distance: 6.053767864070984, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 23.825355146888903, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 24.79825304193105, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(85.73666685032894, 4.105238234975095), expected: [3]NearestResult{{Distance: 30.89312534471159, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 58.167332026144976, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 67.74778455143482, Index: 6, Point: core.NewPdfPoint(70, 70)}}},
		{query: core.NewPdfPoint(70.07227672522, 85.08222783289746), expected: [3]NearestResult{{Distance: 11.281834876246238, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 15.08240101338097, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 21.573202301879537, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(1.845180225020826, 95.15703279645201), expected: [3]NearestResult{{Distance: 33.55349551946908, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 53.12722727818502, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 59.82009650377246, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(77.65420415473436, 0.6457892186903513), expected: [3]NearestResult{{Distance: 22.656973124448452, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 56.573784823694346, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 65.73598041703032, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(6.219304325291253, 47.10473454603812), expected: [3]NearestResult{{Distance: 17.517579846405475, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 20.0796360274705, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 37.47002086164296, Index: 3, Point: core.NewPdfPoint(1, 10)}}},
		{query: core.NewPdfPoint(27.466544481255152, 5.23970743977279), expected: [3]NearestResult{{Distance: 26.891232066181203, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 27.857966400610902, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 30.301027437757085, Index: 2, Point: core.NewPdfPoint(10, 30)}}},
		{query: core.NewPdfPoint(40.721106407631716, 91.17894863211862), expected: [3]NearestResult{{Distance: 5.841316495843984, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 19.168047170329135, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 22.285525137752657, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(93.71967627569863, 69.92495107460142), expected: [3]NearestResult{{Distance: 23.71979500259528, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 35.1926580267403, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 43.02007511262245, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(34.80623433831922, 74.52400508825515), expected: [3]NearestResult{{Distance: 15.477207876099586, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 16.200759780375687, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 25.78201598962232, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(54.410107977828574, 53.02534528663537), expected: [3]NearestResult{{Distance: 5.348061936764951, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 22.23767717618116, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 23.04742145882954, Index: 6, Point: core.NewPdfPoint(70, 70)}}},
		{query: core.NewPdfPoint(75.35946634265616, 43.898478778156736), expected: [3]NearestResult{{Distance: 26.083157293642856, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 26.646074562163907, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 39.23306055946167, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(28.68967491589185, 76.26518874260796), expected: [3]NearestResult{{Distance: 15.115066752856489, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 22.346169871210755, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 31.53228935553006, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(23.173766410901152, 5.938703009992585), expected: [3]NearestResult{{Distance: 22.54262739980084, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 27.43162653380815, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 32.20714100768327, Index: 7, Point: core.NewPdfPoint(55, 1)}}},
		{query: core.NewPdfPoint(45.15121045936782, 84.92310113698109), expected: [3]NearestResult{{Distance: 11.3499769099193, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 11.518518796501736, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 15.64364010155348, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(79.34435328880231, 64.36582210296166), expected: [3]NearestResult{{Distance: 10.91150305693152, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 24.872304329881437, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 30.27355451390367, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(32.71323906831667, 16.467670083297026), expected: [3]NearestResult{{Distance: 24.764179942682542, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 26.43889524827013, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 27.128371322873924, Index: 7, Point: core.NewPdfPoint(55, 1)}}},
		{query: core.NewPdfPoint(7.574147517126595, 2.618810847471864), expected: [3]NearestResult{{Distance: 9.884400279346279, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 27.488439018525362, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 41.24334658113903, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(20.43354407406556, 85.31844765600553), expected: [3]NearestResult{{Distance: 15.300280082134142, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 32.261100258698846, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 39.922303540860256, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(22.34600206444981, 74.86048890604329), expected: [3]NearestResult{{Distance: 19.731407955768056, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 28.654337560583244, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 34.96136999332652, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(63.55251564207386, 87.87501150557617), expected: [3]NearestResult{{Distance: 8.639222974326827, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 17.981422919591974, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 19.002265414160036, Index: 6, Point: core.NewPdfPoint(70, 70)}}},
		{query: core.NewPdfPoint(66.56676934427253, 22.618060969910626), expected: [3]NearestResult{{Distance: 24.51796714987554, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 32.003569043996634, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 45.05472359437686, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(37.37365602876168, 82.29816768422216), expected: [3]NearestResult{{Distance: 8.059309149253208, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 15.457700397197755, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 22.74275744516473, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(49.32206363274456, 34.6288799996447), expected: [3]NearestResult{{Distance: 15.386062777181502, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 24.908065147929335, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 34.104846157417356, Index: 7, Point: core.NewPdfPoint(55, 1)}}},
		{query: core.NewPdfPoint(47.912171617799714, 26.464204062460027), expected: [3]NearestResult{{Distance: 23.628218675284096, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 26.432234103649463, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 26.611752665058166, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(62.98919203656388, 8.505042849821665), expected: [3]NearestResult{{Distance: 10.961425891495823, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 43.48046203362921, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 49.34684425049385, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(22.813112421361513, 91.33114676735282), expected: [3]NearestResult{{Distance: 12.259371132754197, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 32.576172060382156, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 38.874921155541756, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(83.79195528008587, 61.61229643727948), expected: [3]NearestResult{{Distance: 16.142230375755485, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 30.06933285525455, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 35.41952763341757, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(40.57648596422233, 88.72604856926107), expected: [3]NearestResult{{Distance: 5.720152791407797, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 17.23525613908212, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 21.293586384898983, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(23.963796491024425, 39.78301198290163), expected: [3]NearestResult{{Distance: 1.0586791353274034, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 17.049778177452716, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 27.969103262391588, Index: 5, Point: core.NewPdfPoint(50, 50)}}},
		{query: core.NewPdfPoint(66.61803420364912, 81.92497907966498), expected: [3]NearestResult{{Distance: 6.892308842312381, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 12.395274046915407, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 17.084446951544887, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(99.22597668831827, 95.16802163222447), expected: [3]NearestResult{{Distance: 38.569249749849185, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 42.056463562550256, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 52.27326203806389, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(54.91457415134011, 20.76613744667287), expected: [3]NearestResult{{Distance: 19.766322043728387, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 29.64408471981961, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 35.56435315560098, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(37.99829204035451, 49.05561293919847), expected: [3]NearestResult{{Distance: 12.038806455343792, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 15.84170829396003, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 29.019917812230187, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(44.13914280799433, 96.26313880094763), expected: [3]NearestResult{{Distance: 11.07929776226139, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 22.34261473233293, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 22.716876425338192, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(68.45663600770705, 53.9590240469649), expected: [3]NearestResult{{Distance: 16.115051409739802, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 18.876474356336637, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 27.339656357785035, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(12.707495686246627, 54.33094645334382), expected: [3]NearestResult{{Distance: 18.880722670286037, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 24.4811251417603, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 37.54314817876952, Index: 5, Point: core.NewPdfPoint(50, 50)}}},
		{query: core.NewPdfPoint(58.52082047274562, 57.34835769135233), expected: [3]NearestResult{{Distance: 11.25178840401906, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 17.08319688246026, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 19.18705857539686, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(99.15392575500378, 16.024553135770216), expected: [3]NearestResult{{Distance: 46.64017963631754, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 59.75315394816014, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 61.34574354526617, Index: 6, Point: core.NewPdfPoint(70, 70)}}},
		{query: core.NewPdfPoint(9.14386712267613, 18.48796945342841), expected: [3]NearestResult{{Distance: 11.543821326096149, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 11.76300119672437, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 26.724228858097668, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(20.364409084676716, 42.516138522943194), expected: [3]NearestResult{{Distance: 5.274434206705637, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 16.250375355665845, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 30.56593580291902, Index: 5, Point: core.NewPdfPoint(50, 50)}}},
		{query: core.NewPdfPoint(71.41826096965676, 51.52054502261746), expected: [3]NearestResult{{Distance: 18.533799406467093, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 21.472167103721244, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 30.68315562943202, Index: 8, Point: core.NewPdfPoint(60, 80)}}},
		{query: core.NewPdfPoint(66.76607506217822, 41.47325322334792), expected: [3]NearestResult{{Distance: 18.809749694872096, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 28.7094679881515, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 37.048776933820484, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(75.22106099331124, 78.49278018746351), expected: [3]NearestResult{{Distance: 9.96929251293435, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 15.295502911817039, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 24.471602094665585, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(20.23368754320779, 18.77886105160098), expected: [3]NearestResult{{Distance: 15.186912788031798, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 21.142448719887277, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 21.749815463654638, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(26.22566842248103, 2.1646330467549713), expected: [3]NearestResult{{Distance: 26.414528628256093, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 28.79789103157728, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 32.21924842664868, Index: 2, Point: core.NewPdfPoint(10, 30)}}},
		{query: core.NewPdfPoint(73.433313570979, 20.353343954336044), expected: [3]NearestResult{{Distance: 26.727120522437016, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 37.78947471989696, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 49.765229815536166, Index: 6, Point: core.NewPdfPoint(70, 70)}}},
		{query: core.NewPdfPoint(56.567534506738085, 29.559840241162306), expected: [3]NearestResult{{Distance: 21.46934187309903, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 28.602825717584764, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 33.24915293092673, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(58.883110628279766, 29.544477814049795), expected: [3]NearestResult{{Distance: 22.301077156365295, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 28.807390750087734, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 35.45959855989831, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(10.578350777330447, 76.09626135135599), expected: [3]NearestResult{{Distance: 28.102151148353634, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 38.87060650217747, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 40.43651214967753, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(52.351579782397785, 29.629467324272664), expected: [3]NearestResult{{Distance: 20.505816954363393, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 28.75170480024794, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 29.251613025117088, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(5.635680496099216, 31.218456664234452), expected: [3]NearestResult{{Distance: 4.531216323984779, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 21.262463949577448, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 21.718941891213532, Index: 3, Point: core.NewPdfPoint(1, 10)}}},
		{query: core.NewPdfPoint(70.54093686797465, 63.05914799333128), expected: [3]NearestResult{{Distance: 6.961899114006998, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 19.952539105749896, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 22.900483844743007, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(16.722031960407623, 61.0191629618462), expected: [3]NearestResult{{Distance: 22.59048398067558, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 31.739158535322147, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 34.263289846253635, Index: 4, Point: core.NewPdfPoint(35, 90)}}},
		{query: core.NewPdfPoint(2.44084138208438, 6.372760538234356), expected: [3]NearestResult{{Distance: 3.902933512284925, Index: 3, Point: core.NewPdfPoint(1, 10)}, {Distance: 24.807001503495414, Index: 2, Point: core.NewPdfPoint(10, 30)}, {Distance: 40.49329415307187, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(97.04197763063725, 76.19402759440892), expected: [3]NearestResult{{Distance: 27.742287793478468, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 37.236991456624835, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 46.05745765927936, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(3.1014335951180305, 93.45996261954107), expected: [3]NearestResult{{Distance: 32.08566471206863, Index: 4, Point: core.NewPdfPoint(35, 90)}, {Distance: 51.3326687749404, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 57.77122825309978, Index: 1, Point: core.NewPdfPoint(25, 40)}}},
		{query: core.NewPdfPoint(94.23480523412653, 20.79954962662124), expected: [3]NearestResult{{Distance: 43.94760638734355, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 53.00362531100362, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 54.84532889571684, Index: 6, Point: core.NewPdfPoint(70, 70)}}},
		{query: core.NewPdfPoint(92.12116656208707, 68.7190578812228), expected: [3]NearestResult{{Distance: 22.158222464341684, Index: 6, Point: core.NewPdfPoint(70, 70)}, {Distance: 34.04451492379561, Index: 8, Point: core.NewPdfPoint(60, 80)}, {Distance: 41.59808376988461, Index: 0, Point: core.NewPdfPoint(51, 75)}}},
		{query: core.NewPdfPoint(37.35531423574498, 8.63839212418782), expected: [3]NearestResult{{Distance: 19.227063477352917, Index: 7, Point: core.NewPdfPoint(55, 1)}, {Distance: 33.7076287866739, Index: 1, Point: core.NewPdfPoint(25, 40)}, {Distance: 34.7078018315236, Index: 2, Point: core.NewPdfPoint(10, 30)}}},
		{query: core.NewPdfPoint(51.71954908303077, 51.41808119631073), expected: [3]NearestResult{{Distance: 2.228856955545159, Index: 5, Point: core.NewPdfPoint(50, 50)}, {Distance: 23.592893958704682, Index: 0, Point: core.NewPdfPoint(51, 75)}, {Distance: 26.066503259060696, Index: 6, Point: core.NewPdfPoint(70, 70)}}},
	}
	for _, tc := range tests {
		results := tree.FindNearestNeighbours(tc.query, 3, Euclidean)
		if len(results) != 3 { t.Errorf("expected 3 results, got %d", len(results)); continue }
		for i := 0; i < 3; i++ {
			r := results[i]; e := tc.expected[i]
			assertNearKdTree(t, r.Distance, e.Distance, 0.000001)
			if r.Index != e.Index { t.Errorf("[%d] idx: got %d, want %d", i, r.Index, e.Index) }
			if math.Abs(r.Point.X-e.Point.X)>0.001||math.Abs(r.Point.Y-e.Point.Y)>0.001 {
				t.Errorf("[%d] pt: got %v, want %v", i, r.Point, e.Point)
			}
		}
	}
}

// 100 cases from C# DataTree2
func TestNearestTree2Generated(t *testing.T) {
	tree, err := NewKdTree(tree2Points)
	if err != nil { t.Fatal(err) }

	tests := []struct{ query core.PdfPoint; expectDist float64; expectIdx int; expectPt core.PdfPoint }{
		{query: core.NewPdfPoint(87.37822238977932, 47.13424664374255), expectDist: 3.7405807168027154, expectIdx: 77, expectPt: core.NewPdfPoint(84.64548329605812, 44.580018552016185)},
		{query: core.NewPdfPoint(74.5026032549824, 69.93356102115371), expectDist: 7.816974165006036, expectIdx: 73, expectPt: core.NewPdfPoint(82.04208148460329, 71.99786553392418)},
		{query: core.NewPdfPoint(54.01419184925133, 71.12439211149801), expectDist: 3.039056340830242, expectIdx: 71, expectPt: core.NewPdfPoint(52.97863363490906, 68.26721122337617)},
		{query: core.NewPdfPoint(53.56967860131475, 54.19173302307295), expectDist: 10.943391067925589, expectIdx: 60, expectPt: core.NewPdfPoint(43.017566199578106, 57.09185921458071)},
		{query: core.NewPdfPoint(33.62320720050764, 38.388295379543614), expectDist: 3.899048259891558, expectIdx: 56, expectPt: core.NewPdfPoint(29.725739779656745, 38.277277199189065)},
		{query: core.NewPdfPoint(67.71643075209577, 98.50616915055454), expectDist: 6.996934624874407, expectIdx: 75, expectPt: core.NewPdfPoint(62.777729239709636, 93.54972672515291)},
		{query: core.NewPdfPoint(62.853776192101186, 77.58809088626217), expectDist: 7.233120102743402, expectIdx: 32, expectPt: core.NewPdfPoint(55.68259729639997, 78.53266468780453)},
		{query: core.NewPdfPoint(39.1710541173474, 6.731437418177089), expectDist: 3.7156412263891623, expectIdx: 13, expectPt: core.NewPdfPoint(35.803655214742406, 5.16089154047169)},
		{query: core.NewPdfPoint(9.500054562360727, 40.24933325499364), expectDist: 7.832014004033237, expectIdx: 42, expectPt: core.NewPdfPoint(13.34483441571791, 47.0726836958395)},
		{query: core.NewPdfPoint(4.158440156246623, 15.800929441459854), expectDist: 7.657460730577243, expectIdx: 26, expectPt: core.NewPdfPoint(11.228694290441743, 12.860145308234605)},
		{query: core.NewPdfPoint(70.38603008213099, 37.92398478562566), expectDist: 3.5662403566120062, expectIdx: 87, expectPt: core.NewPdfPoint(67.0752294529232, 39.24937886844233)},
		{query: core.NewPdfPoint(36.49654657133572, 22.300307481532634), expectDist: 5.24842528186342, expectIdx: 2, expectPt: core.NewPdfPoint(35.583597244228336, 27.46872033901967)},
		{query: core.NewPdfPoint(44.448096442616745, 76.71337511467348), expectDist: 10.948730989453496, expectIdx: 86, expectPt: core.NewPdfPoint(52.44257355334966, 84.19422038812934)},
		{query: core.NewPdfPoint(18.520340417793502, 44.31476289468993), expectDist: 1.378368419246679, expectIdx: 7, expectPt: core.NewPdfPoint(19.15551535520754, 45.53805856543256)},
		{query: core.NewPdfPoint(68.27767536726084, 94.51997472754229), expectDist: 4.195644188434987, expectIdx: 46, expectPt: core.NewPdfPoint(67.53711275495073, 90.39020510040851)},
		{query: core.NewPdfPoint(80.9852186456766, 37.41169774270004), expectDist: 5.032841375488803, expectIdx: 72, expectPt: core.NewPdfPoint(84.97188896260866, 40.483497230772926)},
		{query: core.NewPdfPoint(57.035762190811745, 39.26879178691382), expectDist: 8.758706221403711, expectIdx: 47, expectPt: core.NewPdfPoint(62.56803790975437, 46.05914229645705)},
		{query: core.NewPdfPoint(12.085443471887913, 97.92593979377935), expectDist: 2.1562961875551574, expectIdx: 18, expectPt: core.NewPdfPoint(14.158071542303697, 98.52077265281578)},
		{query: core.NewPdfPoint(77.27841905993205, 42.15685387387795), expectDist: 4.125216461895987, expectIdx: 74, expectPt: core.NewPdfPoint(75.03458936225277, 45.61844875850991)},
		{query: core.NewPdfPoint(85.87212254856951, 56.063370792337665), expectDist: 3.032017326786329, expectIdx: 29, expectPt: core.NewPdfPoint(86.97060905691454, 53.23733886383501)},
		{query: core.NewPdfPoint(12.121061998950344, 88.5888682638655), expectDist: 10.138645504748764, expectIdx: 18, expectPt: core.NewPdfPoint(14.158071542303697, 98.52077265281578)},
		{query: core.NewPdfPoint(87.40498458296966, 51.58451776586178), expectDist: 1.7089469504759578, expectIdx: 29, expectPt: core.NewPdfPoint(86.97060905691454, 53.23733886383501)},
		{query: core.NewPdfPoint(55.13592516333022, 41.86805042798852), expectDist: 8.532382488232894, expectIdx: 47, expectPt: core.NewPdfPoint(62.56803790975437, 46.05914229645705)},
		{query: core.NewPdfPoint(46.67110223109417, 55.89363812932652), expectDist: 3.8450044606910225, expectIdx: 60, expectPt: core.NewPdfPoint(43.017566199578106, 57.09185921458071)},
		{query: core.NewPdfPoint(16.092850368369838, 39.120936314989095), expectDist: 7.110511570818105, expectIdx: 7, expectPt: core.NewPdfPoint(19.15551535520754, 45.53805856543256)},
		{query: core.NewPdfPoint(60.44407078486588, 51.23153354793123), expectDist: 5.591499584720878, expectIdx: 47, expectPt: core.NewPdfPoint(62.56803790975437, 46.05914229645705)},
		{query: core.NewPdfPoint(16.3698062506329, 73.24470440615686), expectDist: 5.656289708761183, expectIdx: 68, expectPt: core.NewPdfPoint(10.747687004702266, 72.62390777764313)},
		{query: core.NewPdfPoint(23.603385764313224, 35.368200010880244), expectDist: 6.7783441028566624, expectIdx: 56, expectPt: core.NewPdfPoint(29.725739779656745, 38.277277199189065)},
		{query: core.NewPdfPoint(55.00466410944499, 87.18171359059261), expectDist: 3.9356605103079088, expectIdx: 86, expectPt: core.NewPdfPoint(52.44257355334966, 84.19422038812934)},
		{query: core.NewPdfPoint(79.35386322733488, 71.45043817963591), expectDist: 2.7433909868872557, expectIdx: 73, expectPt: core.NewPdfPoint(82.04208148460329, 71.99786553392418)},
		{query: core.NewPdfPoint(92.04745240854216, 81.04432394243022), expectDist: 1.193122715963635, expectIdx: 70, expectPt: core.NewPdfPoint(91.00755813085259, 81.62926984298593)},
		{query: core.NewPdfPoint(38.38165092946303, 48.26803735488079), expectDist: 9.967524396930479, expectIdx: 60, expectPt: core.NewPdfPoint(43.017566199578106, 57.09185921458071)},
		{query: core.NewPdfPoint(76.69366761932686, 86.37384253935019), expectDist: 5.9336956035398005, expectIdx: 41, expectPt: core.NewPdfPoint(74.16588992664359, 81.00550169213861)},
		{query: core.NewPdfPoint(26.857001817848914, 72.32337753974622), expectDist: 7.407056290485798, expectIdx: 10, expectPt: core.NewPdfPoint(24.41810676553939, 79.31739454102991)},
		{query: core.NewPdfPoint(67.72139867299835, 92.06400298246639), expectDist: 1.6839123045966824, expectIdx: 46, expectPt: core.NewPdfPoint(67.53711275495073, 90.39020510040851)},
		{query: core.NewPdfPoint(60.07059131774799, 65.14826966545107), expectDist: 1.1330704932109386, expectIdx: 23, expectPt: core.NewPdfPoint(59.31519405414486, 64.30374393065894)},
		{query: core.NewPdfPoint(98.18469148101073, 19.636669855162936), expectDist: 3.018581465663707, expectIdx: 82, expectPt: core.NewPdfPoint(99.66377946738517, 22.268044551668787)},
		{query: core.NewPdfPoint(13.315505500931046, 2.219169876084226), expectDist: 2.3566520929687598, expectIdx: 64, expectPt: core.NewPdfPoint(15.264247773398054, 0.8939481426395335)},
		{query: core.NewPdfPoint(9.499232323084716, 53.073137917110024), expectDist: 3.9093162473638645, expectIdx: 81, expectPt: core.NewPdfPoint(11.201209707090698, 49.55375920684947)},
		{query: core.NewPdfPoint(68.33664969059441, 84.3521765566152), expectDist: 2.4215364431973803, expectIdx: 58, expectPt: core.NewPdfPoint(70.09436264721477, 86.01779473250726)},
		{query: core.NewPdfPoint(68.17553361331473, 88.07665773162144), expectDist: 2.4000172124415697, expectIdx: 46, expectPt: core.NewPdfPoint(67.53711275495073, 90.39020510040851)},
		{query: core.NewPdfPoint(49.166278548969665, 1.397572786085366), expectDist: 6.93584258247533, expectIdx: 88, expectPt: core.NewPdfPoint(43.03278016511083, 4.6357991376672185)},
		{query: core.NewPdfPoint(24.727638210913018, 32.081396611116276), expectDist: 7.960524829000644, expectIdx: 56, expectPt: core.NewPdfPoint(29.725739779656745, 38.277277199189065)},
		{query: core.NewPdfPoint(88.42089809955208, 39.026696182483356), expectDist: 1.7575731329222621, expectIdx: 79, expectPt: core.NewPdfPoint(90.03637785608008, 39.71900281528301)},
		{query: core.NewPdfPoint(51.42556152624328, 15.448878446223556), expectDist: 1.470978335201828, expectIdx: 34, expectPt: core.NewPdfPoint(52.89633110070122, 15.424096971651757)},
		{query: core.NewPdfPoint(19.693229959708034, 2.7893290005286175), expectDist: 2.1650646266448024, expectIdx: 51, expectPt: core.NewPdfPoint(18.13308334914272, 1.2881803099952127)},
		{query: core.NewPdfPoint(9.279840978602271, 29.123265376879715), expectDist: 6.294219198683488, expectIdx: 57, expectPt: core.NewPdfPoint(15.448731011261907, 30.373061794712797)},
		{query: core.NewPdfPoint(78.20191411050928, 81.43682258294183), expectDist: 2.4567740328680956, expectIdx: 27, expectPt: core.NewPdfPoint(80.62371302859762, 81.02375037506395)},
		{query: core.NewPdfPoint(40.537261184519394, 72.9327512887614), expectDist: 10.517715080249477, expectIdx: 84, expectPt: core.NewPdfPoint(42.195134616764996, 62.54652068319705)},
		{query: core.NewPdfPoint(7.386788690671264, 3.4978575051093697), expectDist: 5.549162421342523, expectIdx: 17, expectPt: core.NewPdfPoint(8.254965696884042, 8.97868520394496)},
		{query: core.NewPdfPoint(26.05756081867855, 84.20846377940514), expectDist: 3.579960160958175, expectIdx: 98, expectPt: core.NewPdfPoint(28.000497306010118, 81.20161949340638)},
		{query: core.NewPdfPoint(74.3139313663137, 31.084939962251912), expectDist: 7.718062249062456, expectIdx: 15, expectPt: core.NewPdfPoint(80.11998397458856, 25.99984035749021)},
		{query: core.NewPdfPoint(27.08914433160744, 37.16399756733006), expectDist: 2.8619970467116924, expectIdx: 56, expectPt: core.NewPdfPoint(29.725739779656745, 38.277277199189065)},
		{query: core.NewPdfPoint(70.85333147347784, 5.242611226453819), expectDist: 1.3605803595971302, expectIdx: 76, expectPt: core.NewPdfPoint(70.36073331446983, 3.9743344551363857)},
		{query: core.NewPdfPoint(51.77204458863274, 52.932206910888034), expectDist: 9.69245062675278, expectIdx: 60, expectPt: core.NewPdfPoint(43.017566199578106, 57.09185921458071)},
		{query: core.NewPdfPoint(13.471218812546915, 11.31698816247696), expectDist: 2.7221774756150476, expectIdx: 26, expectPt: core.NewPdfPoint(11.228694290441743, 12.860145308234605)},
		{query: core.NewPdfPoint(23.382642351306295, 64.21491294035737), expectDist: 5.629527326020367, expectIdx: 69, expectPt: core.NewPdfPoint(18.204503055467892, 62.00619128580146)},
		{query: core.NewPdfPoint(21.987233262429573, 11.042886413409592), expectDist: 6.829449766789799, expectIdx: 1, expectPt: core.NewPdfPoint(16.445674398013832, 7.051331644986569)},
		{query: core.NewPdfPoint(44.65334716075204, 56.72809804882697), expectDist: 1.6757391022022616, expectIdx: 60, expectPt: core.NewPdfPoint(43.017566199578106, 57.09185921458071)},
		{query: core.NewPdfPoint(14.413577866791204, 65.20588702038663), expectDist: 1.2324454456029377, expectIdx: 21, expectPt: core.NewPdfPoint(15.143510163282482, 66.19892308162831)},
		{query: core.NewPdfPoint(1.496525343024535, 34.243471455041515), expectDist: 4.83233101936407, expectIdx: 62, expectPt: core.NewPdfPoint(2.4884293540186175, 29.514037079380717)},
		{query: core.NewPdfPoint(60.37475135874306, 86.38587789183474), expectDist: 7.556125515305361, expectIdx: 75, expectPt: core.NewPdfPoint(62.777729239709636, 93.54972672515291)},
		{query: core.NewPdfPoint(87.96474891454915, 60.71664437764912), expectDist: 5.766963730198733, expectIdx: 0, expectPt: core.NewPdfPoint(82.45353838109239, 62.415005093558115)},
		{query: core.NewPdfPoint(93.90443153077325, 3.5656814220870303), expectDist: 4.363063424232462, expectIdx: 54, expectPt: core.NewPdfPoint(97.7460348272766, 5.634112362720456)},
		{query: core.NewPdfPoint(56.99120651646804, 81.00351342714225), expectDist: 2.7959884805796307, expectIdx: 32, expectPt: core.NewPdfPoint(55.68259729639997, 78.53266468780453)},
		{query: core.NewPdfPoint(76.89422069517484, 54.17932624450519), expectDist: 6.085820174967678, expectIdx: 59, expectPt: core.NewPdfPoint(71.14819197800337, 52.17424240997747)},
		{query: core.NewPdfPoint(63.45013931565355, 80.49930589119406), expectDist: 8.012639180927806, expectIdx: 32, expectPt: core.NewPdfPoint(55.68259729639997, 78.53266468780453)},
		{query: core.NewPdfPoint(31.33405794501082, 30.022880285076926), expectDist: 4.958055796946755, expectIdx: 2, expectPt: core.NewPdfPoint(35.583597244228336, 27.46872033901967)},
		{query: core.NewPdfPoint(59.209873769066554, 46.30109565207685), expectDist: 3.366869142407712, expectIdx: 47, expectPt: core.NewPdfPoint(62.56803790975437, 46.05914229645705)},
		{query: core.NewPdfPoint(78.79398195705156, 78.45536064489083), expectDist: 3.153496725896489, expectIdx: 27, expectPt: core.NewPdfPoint(80.62371302859762, 81.02375037506395)},
		{query: core.NewPdfPoint(85.57104507254752, 24.654183310002807), expectDist: 5.053213509947922, expectIdx: 8, expectPt: core.NewPdfPoint(80.59123120177551, 25.51233282470693)},
		{query: core.NewPdfPoint(78.95708936964519, 12.704872047938153), expectDist: 6.31686834491134, expectIdx: 40, expectPt: core.NewPdfPoint(82.61099784649278, 17.85771383958713)},
		{query: core.NewPdfPoint(13.174935980071456, 88.82914034992214), expectDist: 9.741370141218182, expectIdx: 18, expectPt: core.NewPdfPoint(14.158071542303697, 98.52077265281578)},
		{query: core.NewPdfPoint(55.623046798666124, 96.84465577887713), expectDist: 7.4016913841622145, expectIdx: 36, expectPt: core.NewPdfPoint(49.53804959253269, 92.63065004232082)},
		{query: core.NewPdfPoint(8.430144926873329, 97.99469494927773), expectDist: 5.75203451501427, expectIdx: 18, expectPt: core.NewPdfPoint(14.158071542303697, 98.52077265281578)},
		{query: core.NewPdfPoint(43.22581030476383, 49.88363608992273), expectDist: 7.211230562268732, expectIdx: 60, expectPt: core.NewPdfPoint(43.017566199578106, 57.09185921458071)},
		{query: core.NewPdfPoint(52.462469277256375, 74.30081146714733), expectDist: 5.317688044710818, expectIdx: 32, expectPt: core.NewPdfPoint(55.68259729639997, 78.53266468780453)},
		{query: core.NewPdfPoint(2.0752344418159763, 84.50188522953496), expectDist: 14.707065710166793, expectIdx: 68, expectPt: core.NewPdfPoint(10.747687004702266, 72.62390777764313)},
		{query: core.NewPdfPoint(87.72047866413772, 86.93785567990162), expectDist: 1.8361371712469572, expectIdx: 95, expectPt: core.NewPdfPoint(89.04717161055198, 88.20721990563409)},
		{query: core.NewPdfPoint(85.69225991285919, 6.507776800176323), expectDist: 5.40677558683976, expectIdx: 28, expectPt: core.NewPdfPoint(90.87481745321885, 8.048659530077462)},
		{query: core.NewPdfPoint(95.701988137853, 51.80572553493597), expectDist: 2.11554101881059, expectIdx: 6, expectPt: core.NewPdfPoint(97.13960454465617, 50.25370334786503)},
		{query: core.NewPdfPoint(77.4844196684514, 50.073882954946704), expectDist: 5.084541514035221, expectIdx: 74, expectPt: core.NewPdfPoint(75.03458936225277, 45.61844875850991)},
		{query: core.NewPdfPoint(93.83858926824583, 78.54914176983279), expectDist: 4.183530356997088, expectIdx: 70, expectPt: core.NewPdfPoint(91.00755813085259, 81.62926984298593)},
		{query: core.NewPdfPoint(69.60341179436881, 2.118953273473234), expectDist: 2.003989823845886, expectIdx: 76, expectPt: core.NewPdfPoint(70.36073331446983, 3.9743344551363857)},
		{query: core.NewPdfPoint(21.036323077268815, 30.797805649614528), expectDist: 5.603712380054862, expectIdx: 57, expectPt: core.NewPdfPoint(15.448731011261907, 30.373061794712797)},
		{query: core.NewPdfPoint(67.01310772132939, 11.758224733999745), expectDist: 4.235938831569964, expectIdx: 19, expectPt: core.NewPdfPoint(70.71899728005727, 9.70650341683359)},
		{query: core.NewPdfPoint(58.40755345770724, 54.22645743167164), expectDist: 9.16595152690823, expectIdx: 47, expectPt: core.NewPdfPoint(62.56803790975437, 46.05914229645705)},
		{query: core.NewPdfPoint(49.248852518964966, 70.66641193673007), expectDist: 4.434797767182969, expectIdx: 71, expectPt: core.NewPdfPoint(52.97863363490906, 68.26721122337617)},
		{query: core.NewPdfPoint(2.7054302324065804, 31.088830779825305), expectDist: 1.5896743629376124, expectIdx: 62, expectPt: core.NewPdfPoint(2.4884293540186175, 29.514037079380717)},
		{query: core.NewPdfPoint(30.842485871231062, 31.67560039456263), expectDist: 6.3384522442698294, expectIdx: 2, expectPt: core.NewPdfPoint(35.583597244228336, 27.46872033901967)},
		{query: core.NewPdfPoint(55.56502506452341, 65.54856266577112), expectDist: 3.75239531592163, expectIdx: 71, expectPt: core.NewPdfPoint(52.97863363490906, 68.26721122337617)},
		{query: core.NewPdfPoint(1.267016931392484, 79.30364358908852), expectDist: 11.59749867642406, expectIdx: 68, expectPt: core.NewPdfPoint(10.747687004702266, 72.62390777764313)},
		{query: core.NewPdfPoint(85.80159109593838, 72.9274073673635), expectDist: 3.8727200438238047, expectIdx: 73, expectPt: core.NewPdfPoint(82.04208148460329, 71.99786553392418)},
		{query: core.NewPdfPoint(31.687705038346536, 73.64980455511628), expectDist: 8.403892534033446, expectIdx: 98, expectPt: core.NewPdfPoint(28.000497306010118, 81.20161949340638)},
		{query: core.NewPdfPoint(9.155398598534447, 51.101027419495416), expectDist: 2.565030606787517, expectIdx: 81, expectPt: core.NewPdfPoint(11.201209707090698, 49.55375920684947)},
		{query: core.NewPdfPoint(95.78765937402926, 29.508025935639814), expectDist: 4.180615066997134, expectIdx: 90, expectPt: core.NewPdfPoint(93.45285655237589, 26.040139241642603)},
		{query: core.NewPdfPoint(67.92396367960997, 22.903155180355462), expectDist: 2.5948429360078844, expectIdx: 5, expectPt: core.NewPdfPoint(65.40181734760978, 22.29324720161768)},
		{query: core.NewPdfPoint(38.49112754016583, 70.52567748005846), expectDist: 8.79696604588175, expectIdx: 84, expectPt: core.NewPdfPoint(42.195134616764996, 62.54652068319705)},
		{query: core.NewPdfPoint(59.65725942551804, 65.37065649962103), expectDist: 1.120406688708154, expectIdx: 23, expectPt: core.NewPdfPoint(59.31519405414486, 64.30374393065894)},
		{query: core.NewPdfPoint(45.15904139251526, 74.24820871457037), expectDist: 9.844711972778917, expectIdx: 71, expectPt: core.NewPdfPoint(52.97863363490906, 68.26721122337617)},
	}
	for _, tc := range tests {
		point, idx, dist := tree.FindNearestNeighbour(tc.query, Euclidean)
		assertNearKdTree(t, dist, tc.expectDist, 0.000001)
		if idx != tc.expectIdx { t.Errorf("idx: got %d, want %d", idx, tc.expectIdx) }
		if math.Abs(point.X-tc.expectPt.X)>0.001||math.Abs(point.Y-tc.expectPt.Y)>0.001 {
			t.Errorf("pt: got %v, want %v", point, tc.expectPt)
		}
	}
}

// 100 cases from C# DataTreeK2 (k=3)
func TestNearestKTree2Generated(t *testing.T) {
	tree, err := NewKdTree(tree2Points)
	if err != nil { t.Fatal(err) }

	tests := []struct{ query core.PdfPoint; expected [3]NearestResult }{
		{query: core.NewPdfPoint(90.33048545330094, 47.33938084378586), expected: [3]NearestResult{{Distance: 6.319282378964871, Index: 77, Point: core.NewPdfPoint(84.64548329605812, 44.580018552016185)}, {Distance: 6.787833100869266, Index: 29, Point: core.NewPdfPoint(86.97060905691454, 53.23733886383501)}, {Distance: 7.406576703041731, Index: 6, Point: core.NewPdfPoint(97.13960454465617, 50.25370334786503)}}},
		{query: core.NewPdfPoint(94.00419687661137, 1.6606574220328518), expected: [3]NearestResult{{Distance: 5.457993716991005, Index: 54, Point: core.NewPdfPoint(97.7460348272766, 5.634112362720456)}, {Distance: 7.11333863301438, Index: 28, Point: core.NewPdfPoint(90.87481745321885, 8.048659530077462)}, {Distance: 11.581239683136776, Index: 11, Point: core.NewPdfPoint(98.12630738875605, 12.483472098628233)}}},
		{query: core.NewPdfPoint(85.52978116794068, 97.43823054363845), expected: [3]NearestResult{{Distance: 6.495434817422279, Index: 67, Point: core.NewPdfPoint(91.90691599722378, 96.20420265233045)}, {Distance: 9.733700402810387, Index: 16, Point: core.NewPdfPoint(90.26397977468233, 88.93339150438577)}, {Distance: 9.878440814456654, Index: 95, Point: core.NewPdfPoint(89.04717161055198, 88.20721990563409)}}},
		{query: core.NewPdfPoint(22.62162567299718, 79.83427737771372), expected: [3]NearestResult{{Distance: 1.8693614371543983, Index: 10, Point: core.NewPdfPoint(24.41810676553939, 79.31739454102991)}, {Distance: 5.549944549793103, Index: 98, Point: core.NewPdfPoint(28.000497306010118, 81.20161949340638)}, {Distance: 13.114774665647976, Index: 66, Point: core.NewPdfPoint(34.43836768761709, 85.52303143204219)}}},
		{query: core.NewPdfPoint(7.692210312894288, 89.82892292381919), expected: [3]NearestResult{{Distance: 10.833079578284183, Index: 18, Point: core.NewPdfPoint(14.158071542303697, 98.52077265281578)}, {Distance: 17.47422341605857, Index: 68, Point: core.NewPdfPoint(10.747687004702266, 72.62390777764313)}, {Distance: 19.75469162216385, Index: 10, Point: core.NewPdfPoint(24.41810676553939, 79.31739454102991)}}},
		{query: core.NewPdfPoint(34.94017096790447, 92.33233895331428), expected: [3]NearestResult{{Distance: 6.827772363762564, Index: 66, Point: core.NewPdfPoint(34.43836768761709, 85.52303143204219)}, {Distance: 13.116858855258087, Index: 98, Point: core.NewPdfPoint(28.000497306010118, 81.20161949340638)}, {Distance: 14.60092633517482, Index: 36, Point: core.NewPdfPoint(49.53804959253269, 92.63065004232082)}}},
		{query: core.NewPdfPoint(55.63755197233502, 33.336503157598464), expected: [3]NearestResult{{Distance: 11.95859196631113, Index: 96, Point: core.NewPdfPoint(46.898436066324166, 25.173429288152636)}, {Distance: 12.210892065182815, Index: 80, Point: core.NewPdfPoint(63.98092485758994, 24.420561596263646)}, {Distance: 12.87565785976939, Index: 87, Point: core.NewPdfPoint(67.0752294529232, 39.24937886844233)}}},
		{query: core.NewPdfPoint(81.20904836792344, 64.08948679068317), expected: [3]NearestResult{{Distance: 1.9994837794229934, Index: 30, Point: core.NewPdfPoint(79.77968787046065, 62.69132229661655)}, {Distance: 2.0862991987929322, Index: 0, Point: core.NewPdfPoint(82.45353838109239, 62.415005093558115)}, {Distance: 5.102187165794271, Index: 52, Point: core.NewPdfPoint(81.60160120456186, 59.00242318333139)}}},
		{query: core.NewPdfPoint(71.52680243530935, 60.971051107521035), expected: [3]NearestResult{{Distance: 8.430269922710929, Index: 30, Point: core.NewPdfPoint(79.77968787046065, 62.69132229661655)}, {Distance: 8.804952534770395, Index: 59, Point: core.NewPdfPoint(71.14819197800337, 52.17424240997747)}, {Distance: 10.265333221324632, Index: 52, Point: core.NewPdfPoint(81.60160120456186, 59.00242318333139)}}},
		{query: core.NewPdfPoint(29.27291899970018, 12.928347645791094), expected: [3]NearestResult{{Distance: 9.723303049123018, Index: 99, Point: core.NewPdfPoint(24.103817108644833, 21.163820177769633)}, {Distance: 10.14809784435131, Index: 13, Point: core.NewPdfPoint(35.803655214742406, 5.16089154047169)}, {Distance: 10.227363038462595, Index: 38, Point: core.NewPdfPoint(24.636942156010523, 22.044633809864933)}}},
		{query: core.NewPdfPoint(46.78746330316932, 67.20315193481265), expected: [3]NearestResult{{Distance: 6.281943349489286, Index: 71, Point: core.NewPdfPoint(52.97863363490906, 68.26721122337617)}, {Distance: 6.540160347995691, Index: 84, Point: core.NewPdfPoint(42.195134616764996, 62.54652068319705)}, {Distance: 10.791217014122221, Index: 60, Point: core.NewPdfPoint(43.017566199578106, 57.09185921458071)}}},
		{query: core.NewPdfPoint(89.49370609342411, 49.19740999815338), expected: [3]NearestResult{{Distance: 4.763091842008821, Index: 29, Point: core.NewPdfPoint(86.97060905691454, 53.23733886383501)}, {Distance: 6.695189919618434, Index: 77, Point: core.NewPdfPoint(84.64548329605812, 44.580018552016185)}, {Distance: 7.718517912604583, Index: 6, Point: core.NewPdfPoint(97.13960454465617, 50.25370334786503)}}},
		{query: core.NewPdfPoint(50.085173044511734, 76.09967482070661), expected: [3]NearestResult{{Distance: 6.103326793563302, Index: 32, Point: core.NewPdfPoint(55.68259729639997, 78.53266468780453)}, {Distance: 8.34982635697827, Index: 71, Point: core.NewPdfPoint(52.97863363490906, 68.26721122337617)}, {Distance: 8.430836560042678, Index: 86, Point: core.NewPdfPoint(52.44257355334966, 84.19422038812934)}}},
		{query: core.NewPdfPoint(51.41806121874436, 6.185936468302689), expected: [3]NearestResult{{Distance: 4.63313326645828, Index: 43, Point: core.NewPdfPoint(47.945049322341525, 9.252548969452246)}, {Distance: 5.775831321961968, Index: 83, Point: core.NewPdfPoint(57.19140901930086, 6.016576988006017)}, {Distance: 6.717112422982287, Index: 78, Point: core.NewPdfPoint(44.866927314279856, 7.669932380951572)}}},
		{query: core.NewPdfPoint(44.84550779030454, 6.558735867057342), expected: [3]NearestResult{{Distance: 1.1114029370565919, Index: 78, Point: core.NewPdfPoint(44.866927314279856, 7.669932380951572)}, {Distance: 2.6426628820903373, Index: 88, Point: core.NewPdfPoint(43.03278016511083, 4.6357991376672185)}, {Distance: 4.106554119874308, Index: 43, Point: core.NewPdfPoint(47.945049322341525, 9.252548969452246)}}},
		{query: core.NewPdfPoint(61.96609728269268, 15.51346888405012), expected: [3]NearestResult{{Distance: 7.600629342352724, Index: 5, Point: core.NewPdfPoint(65.40181734760978, 22.29324720161768)}, {Distance: 9.070206499012121, Index: 34, Point: core.NewPdfPoint(52.89633110070122, 15.424096971651757)}, {Distance: 9.132131774155129, Index: 80, Point: core.NewPdfPoint(63.98092485758994, 24.420561596263646)}}},
		{query: core.NewPdfPoint(69.5760144871494, 76.4622356528522), expected: [3]NearestResult{{Distance: 6.458190369895054, Index: 41, Point: core.NewPdfPoint(74.16588992664359, 81.00550169213861)}, {Distance: 9.569607836260653, Index: 58, Point: core.NewPdfPoint(70.09436264721477, 86.01779473250726)}, {Distance: 11.952366277171725, Index: 27, Point: core.NewPdfPoint(80.62371302859762, 81.02375037506395)}}},
		{query: core.NewPdfPoint(82.98110145635202, 52.72268932669417), expected: [3]NearestResult{{Distance: 4.022565728614683, Index: 29, Point: core.NewPdfPoint(86.97060905691454, 53.23733886383501)}, {Distance: 6.429469515822035, Index: 52, Point: core.NewPdfPoint(81.60160120456186, 59.00242318333139)}, {Distance: 8.311032081103916, Index: 77, Point: core.NewPdfPoint(84.64548329605812, 44.580018552016185)}}},
		{query: core.NewPdfPoint(50.69508238708692, 87.63717989744282), expected: [3]NearestResult{{Distance: 3.861048505126353, Index: 86, Point: core.NewPdfPoint(52.44257355334966, 84.19422038812934)}, {Distance: 5.125765208772435, Index: 36, Point: core.NewPdfPoint(49.53804959253269, 92.63065004232082)}, {Distance: 10.381112761797503, Index: 32, Point: core.NewPdfPoint(55.68259729639997, 78.53266468780453)}}},
		{query: core.NewPdfPoint(45.394375588571556, 61.05494327358284), expected: [3]NearestResult{{Distance: 3.529864864914409, Index: 84, Point: core.NewPdfPoint(42.195134616764996, 62.54652068319705)}, {Distance: 4.6211749729180625, Index: 60, Point: core.NewPdfPoint(43.017566199578106, 57.09185921458071)}, {Distance: 10.4660297674453, Index: 71, Point: core.NewPdfPoint(52.97863363490906, 68.26721122337617)}}},
		{query: core.NewPdfPoint(21.560319619569214, 46.112818457615276), expected: [3]NearestResult{{Distance: 2.4725356384800556, Index: 7, Point: core.NewPdfPoint(19.15551535520754, 45.53805856543256)}, {Distance: 8.271368593543086, Index: 42, Point: core.NewPdfPoint(13.34483441571791, 47.0726836958395)}, {Distance: 10.915641594452936, Index: 81, Point: core.NewPdfPoint(11.201209707090698, 49.55375920684947)}}},
		{query: core.NewPdfPoint(13.31631815615828, 49.618965383452654), expected: [3]NearestResult{{Distance: 2.116113323237742, Index: 81, Point: core.NewPdfPoint(11.201209707090698, 49.55375920684947)}, {Distance: 2.5464413619271418, Index: 42, Point: core.NewPdfPoint(13.34483441571791, 47.0726836958395)}, {Distance: 7.1239051360014365, Index: 7, Point: core.NewPdfPoint(19.15551535520754, 45.53805856543256)}}},
		{query: core.NewPdfPoint(78.68697517856405, 37.95356900539214), expected: [3]NearestResult{{Distance: 6.75635336262533, Index: 33, Point: core.NewPdfPoint(84.80993603153726, 40.80973096746926)}, {Distance: 6.775003918703829, Index: 72, Point: core.NewPdfPoint(84.97188896260866, 40.483497230772926)}, {Distance: 8.49060090811873, Index: 74, Point: core.NewPdfPoint(75.03458936225277, 45.61844875850991)}}},
		{query: core.NewPdfPoint(40.25217869763216, 47.69662857204957), expected: [3]NearestResult{{Distance: 9.793759587731104, Index: 60, Point: core.NewPdfPoint(43.017566199578106, 57.09185921458071)}, {Distance: 11.423458353755956, Index: 49, Point: core.NewPdfPoint(37.970724110053844, 36.503310666941424)}, {Distance: 11.98058569840965, Index: 35, Point: core.NewPdfPoint(31.112041569584935, 55.44209995586027)}}},
		{query: core.NewPdfPoint(54.1283188372334, 86.53269307585674), expected: [3]NearestResult{{Distance: 2.8827403062682038, Index: 86, Point: core.NewPdfPoint(52.44257355334966, 84.19422038812934)}, {Distance: 7.632538955268626, Index: 36, Point: core.NewPdfPoint(49.53804959253269, 92.63065004232082)}, {Distance: 8.149615680403004, Index: 32, Point: core.NewPdfPoint(55.68259729639997, 78.53266468780453)}}},
		{query: core.NewPdfPoint(81.21021662298843, 89.2711165353666), expected: [3]NearestResult{{Distance: 4.055682992421869, Index: 44, Point: core.NewPdfPoint(82.78544384909999, 85.53384162961818)}, {Distance: 7.031757366483974, Index: 45, Point: core.NewPdfPoint(86.41791572882465, 84.5461317010868)}, {Distance: 7.908839327983186, Index: 95, Point: core.NewPdfPoint(89.04717161055198, 88.20721990563409)}}},
		{query: core.NewPdfPoint(41.04671360406829, 96.94952583007202), expected: [3]NearestResult{{Distance: 9.526566797068826, Index: 36, Point: core.NewPdfPoint(49.53804959253269, 92.63065004232082)}, {Distance: 13.199810982725637, Index: 66, Point: core.NewPdfPoint(34.43836768761709, 85.52303143204219)}, {Distance: 17.104485987625814, Index: 86, Point: core.NewPdfPoint(52.44257355334966, 84.19422038812934)}}},
		{query: core.NewPdfPoint(20.268880198632655, 20.520060156880948), expected: [3]NearestResult{{Distance: 2.0060215634503926, Index: 50, Point: core.NewPdfPoint(18.3788032374675, 19.84795034896214)}, {Distance: 3.888594613516348, Index: 99, Point: core.NewPdfPoint(24.103817108644833, 21.163820177769633)}, {Distance: 4.626477070824485, Index: 38, Point: core.NewPdfPoint(24.636942156010523, 22.044633809864933)}}},
		{query: core.NewPdfPoint(38.68584315271796, 74.86521565365805), expected: [3]NearestResult{{Distance: 11.473015514354003, Index: 66, Point: core.NewPdfPoint(34.43836768761709, 85.52303143204219)}, {Distance: 12.422826952193041, Index: 98, Point: core.NewPdfPoint(28.000497306010118, 81.20161949340638)}, {Distance: 12.808800582212788, Index: 84, Point: core.NewPdfPoint(42.195134616764996, 62.54652068319705)}}},
		{query: core.NewPdfPoint(67.4231014401171, 82.7767352052507), expected: [3]NearestResult{{Distance: 4.200012297096924, Index: 58, Point: core.NewPdfPoint(70.09436264721477, 86.01779473250726)}, {Distance: 6.971546796228628, Index: 41, Point: core.NewPdfPoint(74.16588992664359, 81.00550169213861)}, {Distance: 7.61432350405368, Index: 46, Point: core.NewPdfPoint(67.53711275495073, 90.39020510040851)}}},
		{query: core.NewPdfPoint(14.38235206959435, 13.07129438479172), expected: [3]NearestResult{{Distance: 3.160718481696952, Index: 26, Point: core.NewPdfPoint(11.228694290441743, 12.860145308234605)}, {Distance: 6.363744999573501, Index: 1, Point: core.NewPdfPoint(16.445674398013832, 7.051331644986569)}, {Distance: 7.3684675250439415, Index: 17, Point: core.NewPdfPoint(8.254965696884042, 8.97868520394496)}}},
		{query: core.NewPdfPoint(52.14368788446655, 22.1218586105557), expected: [3]NearestResult{{Distance: 3.2956133368175657, Index: 85, Point: core.NewPdfPoint(49.61310425442538, 20.010647167526496)}, {Distance: 4.200265094831746, Index: 89, Point: core.NewPdfPoint(49.42598643546282, 18.91930873569896)}, {Distance: 6.0683399901534365, Index: 96, Point: core.NewPdfPoint(46.898436066324166, 25.173429288152636)}}},
		{query: core.NewPdfPoint(86.29153320726485, 95.11424407263593), expected: [3]NearestResult{{Distance: 5.720186498989052, Index: 67, Point: core.NewPdfPoint(91.90691599722378, 96.20420265233045)}, {Distance: 7.347330821559044, Index: 16, Point: core.NewPdfPoint(90.26397977468233, 88.93339150438577)}, {Distance: 7.43643233366769, Index: 95, Point: core.NewPdfPoint(89.04717161055198, 88.20721990563409)}}},
		{query: core.NewPdfPoint(81.89614459928313, 66.58499042641813), expected: [3]NearestResult{{Distance: 3.6799201803348827, Index: 55, Point: core.NewPdfPoint(84.9666796464414, 68.61319827025679)}, {Distance: 4.207073270608355, Index: 0, Point: core.NewPdfPoint(82.45353838109239, 62.415005093558115)}, {Distance: 4.431708540733531, Index: 30, Point: core.NewPdfPoint(79.77968787046065, 62.69132229661655)}}},
		{query: core.NewPdfPoint(46.050742000044394, 44.44807688975785), expected: [3]NearestResult{{Distance: 11.331637103162638, Index: 49, Point: core.NewPdfPoint(37.970724110053844, 36.503310666941424)}, {Distance: 13.002514638101248, Index: 60, Point: core.NewPdfPoint(43.017566199578106, 57.09185921458071)}, {Distance: 16.595680037696102, Index: 47, Point: core.NewPdfPoint(62.56803790975437, 46.05914229645705)}}},
		{query: core.NewPdfPoint(2.4117277982407592, 98.30730212163033), expected: [3]NearestResult{{Distance: 11.748283322314418, Index: 18, Point: core.NewPdfPoint(14.158071542303697, 98.52077265281578)}, {Distance: 27.00231399196269, Index: 68, Point: core.NewPdfPoint(10.747687004702266, 72.62390777764313)}, {Distance: 29.067117249085932, Index: 10, Point: core.NewPdfPoint(24.41810676553939, 79.31739454102991)}}},
		{query: core.NewPdfPoint(39.07511433748291, 61.28307901688501), expected: [3]NearestResult{{Distance: 3.3661270604813507, Index: 84, Point: core.NewPdfPoint(42.195134616764996, 62.54652068319705)}, {Distance: 5.754063791457791, Index: 60, Point: core.NewPdfPoint(43.017566199578106, 57.09185921458071)}, {Distance: 9.875604502923744, Index: 35, Point: core.NewPdfPoint(31.112041569584935, 55.44209995586027)}}},
		{query: core.NewPdfPoint(9.49024018298229, 59.73883570183125), expected: [3]NearestResult{{Distance: 5.795324777487998, Index: 9, Point: core.NewPdfPoint(5.021252133213605, 56.04912906205743)}, {Distance: 7.206582635751992, Index: 22, Point: core.NewPdfPoint(16.49098129875618, 61.44894011847811)}, {Distance: 8.338201846917082, Index: 97, Point: core.NewPdfPoint(8.420107657979937, 68.0080815210619)}}},
		{query: core.NewPdfPoint(12.96964790803895, 56.98374692967552), expected: [3]NearestResult{{Distance: 5.686628092455948, Index: 22, Point: core.NewPdfPoint(16.49098129875618, 61.44894011847811)}, {Distance: 7.25456102910056, Index: 69, Point: core.NewPdfPoint(18.204503055467892, 62.00619128580146)}, {Distance: 7.637544843201832, Index: 81, Point: core.NewPdfPoint(11.201209707090698, 49.55375920684947)}}},
		{query: core.NewPdfPoint(13.307937277045156, 52.91717446722356), expected: [3]NearestResult{{Distance: 3.968735726616667, Index: 81, Point: core.NewPdfPoint(11.201209707090698, 49.55375920684947)}, {Distance: 5.844607238783087, Index: 42, Point: core.NewPdfPoint(13.34483441571791, 47.0726836958395)}, {Distance: 8.858797336947017, Index: 9, Point: core.NewPdfPoint(5.021252133213605, 56.04912906205743)}}},
		{query: core.NewPdfPoint(26.36175597577918, 31.381877739024254), expected: [3]NearestResult{{Distance: 7.672217459639145, Index: 56, Point: core.NewPdfPoint(29.725739779656745, 38.277277199189065)}, {Distance: 9.495214947829645, Index: 38, Point: core.NewPdfPoint(24.636942156010523, 22.044633809864933)}, {Distance: 10.01774212173994, Index: 2, Point: core.NewPdfPoint(35.583597244228336, 27.46872033901967)}}},
		{query: core.NewPdfPoint(71.51791012620515, 60.383581685691176), expected: [3]NearestResult{{Distance: 8.21766042452993, Index: 59, Point: core.NewPdfPoint(71.14819197800337, 52.17424240997747)}, {Distance: 8.578032304834814, Index: 30, Point: core.NewPdfPoint(79.77968787046065, 62.69132229661655)}, {Distance: 10.17783987751681, Index: 52, Point: core.NewPdfPoint(81.60160120456186, 59.00242318333139)}}},
		{query: core.NewPdfPoint(13.820958137480366, 45.57166856903693), expected: [3]NearestResult{{Distance: 1.5747190890171712, Index: 42, Point: core.NewPdfPoint(13.34483441571791, 47.0726836958395)}, {Distance: 4.766563509099969, Index: 81, Point: core.NewPdfPoint(11.201209707090698, 49.55375920684947)}, {Distance: 5.334663095411686, Index: 7, Point: core.NewPdfPoint(19.15551535520754, 45.53805856543256)}}},
		{query: core.NewPdfPoint(62.60837807772033, 59.51480363049938), expected: [3]NearestResult{{Distance: 5.8119712844804265, Index: 23, Point: core.NewPdfPoint(59.31519405414486, 64.30374393065894)}, {Distance: 10.500007541726367, Index: 20, Point: core.NewPdfPoint(58.60508497220939, 69.22169822597527)}, {Distance: 10.892398267625163, Index: 24, Point: core.NewPdfPoint(58.880198617023936, 69.74930498275906)}}},
		{query: core.NewPdfPoint(54.246412695504056, 33.090909319980845), expected: [3]NearestResult{{Distance: 10.801817004438293, Index: 96, Point: core.NewPdfPoint(46.898436066324166, 25.173429288152636)}, {Distance: 13.035937123351204, Index: 80, Point: core.NewPdfPoint(63.98092485758994, 24.420561596263646)}, {Distance: 13.876628015735688, Index: 85, Point: core.NewPdfPoint(49.61310425442538, 20.010647167526496)}}},
		{query: core.NewPdfPoint(37.59119855061173, 20.84553738123185), expected: [3]NearestResult{{Distance: 6.920766973227989, Index: 2, Point: core.NewPdfPoint(35.583597244228336, 27.46872033901967)}, {Distance: 7.1217149520573555, Index: 53, Point: core.NewPdfPoint(37.57869607184197, 27.967241358959505)}, {Distance: 10.264273892091708, Index: 96, Point: core.NewPdfPoint(46.898436066324166, 25.173429288152636)}}},
		{query: core.NewPdfPoint(60.2394597901307, 52.39882244517341), expected: [3]NearestResult{{Distance: 6.753800444728761, Index: 47, Point: core.NewPdfPoint(62.56803790975437, 46.05914229645705)}, {Distance: 8.316180599171837, Index: 3, Point: core.NewPdfPoint(63.8245554804822, 44.89509346225664)}, {Distance: 10.911043677803592, Index: 59, Point: core.NewPdfPoint(71.14819197800337, 52.17424240997747)}}},
		{query: core.NewPdfPoint(69.05917872300058, 14.344285937634815), expected: [3]NearestResult{{Distance: 4.92585265234346, Index: 19, Point: core.NewPdfPoint(70.71899728005727, 9.70650341683359)}, {Distance: 6.728125898818047, Index: 48, Point: core.NewPdfPoint(71.56130305322326, 8.098723357173876)}, {Distance: 8.314145698967357, Index: 92, Point: core.NewPdfPoint(69.38884149074713, 6.0366785105100185)}}},
		{query: core.NewPdfPoint(51.331730724820865, 96.31696159244558), expected: [3]NearestResult{{Distance: 4.099534711270266, Index: 36, Point: core.NewPdfPoint(49.53804959253269, 92.63065004232082)}, {Distance: 11.775757759634676, Index: 75, Point: core.NewPdfPoint(62.777729239709636, 93.54972672515291)}, {Distance: 12.173529730383077, Index: 86, Point: core.NewPdfPoint(52.44257355334966, 84.19422038812934)}}},
		{query: core.NewPdfPoint(80.9218188849778, 57.661185428982684), expected: [3]NearestResult{{Distance: 1.5036697495492302, Index: 52, Point: core.NewPdfPoint(81.60160120456186, 59.00242318333139)}, {Distance: 4.994493569730729, Index: 0, Point: core.NewPdfPoint(82.45353838109239, 62.415005093558115)}, {Distance: 5.1581721725288885, Index: 30, Point: core.NewPdfPoint(79.77968787046065, 62.69132229661655)}}},
		{query: core.NewPdfPoint(4.784694797817046, 62.959677091870205), expected: [3]NearestResult{{Distance: 6.221142495114544, Index: 97, Point: core.NewPdfPoint(8.420107657979937, 68.0080815210619)}, {Distance: 6.9145956819816385, Index: 9, Point: core.NewPdfPoint(5.021252133213605, 56.04912906205743)}, {Distance: 10.853468125809684, Index: 21, Point: core.NewPdfPoint(15.143510163282482, 66.19892308162831)}}},
		{query: core.NewPdfPoint(82.41590594603949, 28.680736275752395), expected: [3]NearestResult{{Distance: 0.6666342489180029, Index: 12, Point: core.NewPdfPoint(82.48631571672695, 29.34364176375367)}, {Distance: 3.5296544623441446, Index: 15, Point: core.NewPdfPoint(80.11998397458856, 25.99984035749021)}, {Distance: 3.65625742405422, Index: 8, Point: core.NewPdfPoint(80.59123120177551, 25.51233282470693)}}},
		{query: core.NewPdfPoint(58.93828314851418, 41.6374695944503), expected: [3]NearestResult{{Distance: 5.72069131402993, Index: 47, Point: core.NewPdfPoint(62.56803790975437, 46.05914229645705)}, {Distance: 5.8726289314290705, Index: 3, Point: core.NewPdfPoint(63.8245554804822, 44.89509346225664)}, {Distance: 8.480145781558287, Index: 87, Point: core.NewPdfPoint(67.0752294529232, 39.24937886844233)}}},
		{query: core.NewPdfPoint(87.95455074935215, 32.292861392280614), expected: [3]NearestResult{{Distance: 6.212848846488876, Index: 12, Point: core.NewPdfPoint(82.48631571672695, 29.34364176375367)}, {Distance: 7.712430261385789, Index: 79, Point: core.NewPdfPoint(90.03637785608008, 39.71900281528301)}, {Distance: 8.326337790207875, Index: 90, Point: core.NewPdfPoint(93.45285655237589, 26.040139241642603)}}},
		{query: core.NewPdfPoint(79.20389918634515, 25.654064891808016), expected: [3]NearestResult{{Distance: 0.9791690415442634, Index: 15, Point: core.NewPdfPoint(80.11998397458856, 25.99984035749021)}, {Distance: 1.3945530107825985, Index: 8, Point: core.NewPdfPoint(80.59123120177551, 25.51233282470693)}, {Distance: 4.081864905775139, Index: 94, Point: core.NewPdfPoint(80.69475374791655, 21.854199922464968)}}},
		{query: core.NewPdfPoint(63.37032482613486, 52.83441603652857), expected: [3]NearestResult{{Distance: 6.822609365125496, Index: 47, Point: core.NewPdfPoint(62.56803790975437, 46.05914229645705)}, {Distance: 7.805834141801166, Index: 59, Point: core.NewPdfPoint(71.14819197800337, 52.17424240997747)}, {Distance: 7.952305855894424, Index: 3, Point: core.NewPdfPoint(63.8245554804822, 44.89509346225664)}}},
		{query: core.NewPdfPoint(66.0404611511353, 76.6960457458976), expected: [3]NearestResult{{Distance: 9.197499841720573, Index: 41, Point: core.NewPdfPoint(74.16588992664359, 81.00550169213861)}, {Distance: 9.976300255487988, Index: 24, Point: core.NewPdfPoint(58.880198617023936, 69.74930498275906)}, {Distance: 10.165093285812645, Index: 58, Point: core.NewPdfPoint(70.09436264721477, 86.01779473250726)}}},
		{query: core.NewPdfPoint(73.6630405559061, 20.85637324912446), expected: [3]NearestResult{{Distance: 7.10215801603454, Index: 94, Point: core.NewPdfPoint(80.69475374791655, 21.854199922464968)}, {Distance: 8.25514216757696, Index: 15, Point: core.NewPdfPoint(80.11998397458856, 25.99984035749021)}, {Distance: 8.347322037334528, Index: 8, Point: core.NewPdfPoint(80.59123120177551, 25.51233282470693)}}},
		{query: core.NewPdfPoint(97.68962763674809, 43.53575778881722), expected: [3]NearestResult{{Distance: 4.718909708545427, Index: 93, Point: core.NewPdfPoint(99.69679458529599, 47.806517636673476)}, {Distance: 6.740424165893753, Index: 6, Point: core.NewPdfPoint(97.13960454465617, 50.25370334786503)}, {Distance: 8.55218397447652, Index: 79, Point: core.NewPdfPoint(90.03637785608008, 39.71900281528301)}}},
		{query: core.NewPdfPoint(3.9643845233934605, 49.339573111296765), expected: [3]NearestResult{{Distance: 2.995004662567729, Index: 14, Point: core.NewPdfPoint(1.813800375443253, 47.255096977735214)}, {Distance: 6.792283136109303, Index: 9, Point: core.NewPdfPoint(5.021252133213605, 56.04912906205743)}, {Distance: 7.239994089978434, Index: 81, Point: core.NewPdfPoint(11.201209707090698, 49.55375920684947)}}},
		{query: core.NewPdfPoint(70.73913716966032, 10.028007391201522), expected: [3]NearestResult{{Distance: 0.3221341656633137, Index: 19, Point: core.NewPdfPoint(70.71899728005727, 9.70650341683359)}, {Distance: 2.0971632325712752, Index: 48, Point: core.NewPdfPoint(71.56130305322326, 8.098723357173876)}, {Distance: 4.2135501248156215, Index: 92, Point: core.NewPdfPoint(69.38884149074713, 6.0366785105100185)}}},
		{query: core.NewPdfPoint(79.28915979507235, 54.70843490604055), expected: [3]NearestResult{{Distance: 4.877060651440881, Index: 52, Point: core.NewPdfPoint(81.60160120456186, 59.00242318333139)}, {Distance: 7.82104764898212, Index: 29, Point: core.NewPdfPoint(86.97060905691454, 53.23733886383501)}, {Distance: 7.997944041024629, Index: 30, Point: core.NewPdfPoint(79.77968787046065, 62.69132229661655)}}},
		{query: core.NewPdfPoint(87.95765064513492, 76.78785901003045), expected: [3]NearestResult{{Distance: 5.721992181478006, Index: 70, Point: core.NewPdfPoint(91.00755813085259, 81.62926984298593)}, {Distance: 6.098901750707647, Index: 37, Point: core.NewPdfPoint(92.45607820110665, 72.66951564303699)}, {Distance: 6.409897446967801, Index: 65, Point: core.NewPdfPoint(82.27200991589356, 79.74763514660499)}}},
		{query: core.NewPdfPoint(8.309511701427963, 39.13056340529091), expected: [3]NearestResult{{Distance: 9.403815690802748, Index: 42, Point: core.NewPdfPoint(13.34483441571791, 47.0726836958395)}, {Distance: 10.40203400303521, Index: 14, Point: core.NewPdfPoint(1.813800375443253, 47.255096977735214)}, {Distance: 10.816881624275172, Index: 81, Point: core.NewPdfPoint(11.201209707090698, 49.55375920684947)}}},
		{query: core.NewPdfPoint(39.732285598197045, 75.85998020600681), expected: [3]NearestResult{{Distance: 11.018172527281603, Index: 66, Point: core.NewPdfPoint(34.43836768761709, 85.52303143204219)}, {Distance: 12.890615455027152, Index: 98, Point: core.NewPdfPoint(28.000497306010118, 81.20161949340638)}, {Distance: 13.53934377116391, Index: 84, Point: core.NewPdfPoint(42.195134616764996, 62.54652068319705)}}},
		{query: core.NewPdfPoint(42.52668278329558, 59.630604107901206), expected: [3]NearestResult{{Distance: 2.5857672288398494, Index: 60, Point: core.NewPdfPoint(43.017566199578106, 57.09185921458071)}, {Distance: 2.9347050381281665, Index: 84, Point: core.NewPdfPoint(42.195134616764996, 62.54652068319705)}, {Distance: 12.158848673678182, Index: 35, Point: core.NewPdfPoint(31.112041569584935, 55.44209995586027)}}},
		{query: core.NewPdfPoint(11.91848255296697, 43.04840499795828), expected: [3]NearestResult{{Distance: 4.269578278307316, Index: 42, Point: core.NewPdfPoint(13.34483441571791, 47.0726836958395)}, {Distance: 6.544777591222621, Index: 81, Point: core.NewPdfPoint(11.201209707090698, 49.55375920684947)}, {Distance: 7.6533011613775805, Index: 7, Point: core.NewPdfPoint(19.15551535520754, 45.53805856543256)}}},
		{query: core.NewPdfPoint(45.884180750302804, 84.54782622509595), expected: [3]NearestResult{{Distance: 6.567918486628291, Index: 86, Point: core.NewPdfPoint(52.44257355334966, 84.19422038812934)}, {Distance: 8.870332483989317, Index: 36, Point: core.NewPdfPoint(49.53804959253269, 92.63065004232082)}, {Distance: 11.487282614334951, Index: 66, Point: core.NewPdfPoint(34.43836768761709, 85.52303143204219)}}},
		{query: core.NewPdfPoint(55.67904148326467, 33.58226552458229), expected: [3]NearestResult{{Distance: 12.157613184301322, Index: 96, Point: core.NewPdfPoint(46.898436066324166, 25.173429288152636)}, {Distance: 12.363579030000482, Index: 80, Point: core.NewPdfPoint(63.98092485758994, 24.420561596263646)}, {Distance: 12.727500693064066, Index: 87, Point: core.NewPdfPoint(67.0752294529232, 39.24937886844233)}}},
		{query: core.NewPdfPoint(62.05241890629847, 18.690682757200527), expected: [3]NearestResult{{Distance: 4.919038574237582, Index: 5, Point: core.NewPdfPoint(65.40181734760978, 22.29324720161768)}, {Distance: 6.045713085692203, Index: 80, Point: core.NewPdfPoint(63.98092485758994, 24.420561596263646)}, {Distance: 9.721343867910239, Index: 34, Point: core.NewPdfPoint(52.89633110070122, 15.424096971651757)}}},
		{query: core.NewPdfPoint(11.707658570636525, 65.12592360890554), expected: [3]NearestResult{{Distance: 3.5995005257884185, Index: 21, Point: core.NewPdfPoint(15.143510163282482, 66.19892308162831)}, {Distance: 4.372050461043956, Index: 97, Point: core.NewPdfPoint(8.420107657979937, 68.0080815210619)}, {Distance: 6.033273067765282, Index: 22, Point: core.NewPdfPoint(16.49098129875618, 61.44894011847811)}}},
		{query: core.NewPdfPoint(33.90020310829871, 47.7313194477634), expected: [3]NearestResult{{Distance: 8.199389051021392, Index: 35, Point: core.NewPdfPoint(31.112041569584935, 55.44209995586027)}, {Distance: 10.334653304296424, Index: 56, Point: core.NewPdfPoint(29.725739779656745, 38.277277199189065)}, {Distance: 11.943086804002764, Index: 49, Point: core.NewPdfPoint(37.970724110053844, 36.503310666941424)}}},
		{query: core.NewPdfPoint(42.15217638961487, 22.384017040998284), expected: [3]NearestResult{{Distance: 5.505252183445487, Index: 96, Point: core.NewPdfPoint(46.898436066324166, 25.173429288152636)}, {Distance: 7.217278988769222, Index: 53, Point: core.NewPdfPoint(37.57869607184197, 27.967241358959505)}, {Distance: 7.829324949202796, Index: 85, Point: core.NewPdfPoint(49.61310425442538, 20.010647167526496)}}},
		{query: core.NewPdfPoint(59.752335901407626, 57.2509255780393), expected: [3]NearestResult{{Distance: 7.066352645437156, Index: 23, Point: core.NewPdfPoint(59.31519405414486, 64.30374393065894)}, {Distance: 11.540545516643135, Index: 47, Point: core.NewPdfPoint(62.56803790975437, 46.05914229645705)}, {Distance: 12.025621916687646, Index: 20, Point: core.NewPdfPoint(58.60508497220939, 69.22169822597527)}}},
		{query: core.NewPdfPoint(97.9181308204627, 34.375867392622084), expected: [3]NearestResult{{Distance: 0.9254860223406257, Index: 31, Point: core.NewPdfPoint(97.45369858932968, 33.575350634376086)}, {Distance: 9.456375526398574, Index: 90, Point: core.NewPdfPoint(93.45285655237589, 26.040139241642603)}, {Distance: 9.522138727011285, Index: 79, Point: core.NewPdfPoint(90.03637785608008, 39.71900281528301)}}},
		{query: core.NewPdfPoint(68.84376255969894, 78.95185188531323), expected: [3]NearestResult{{Distance: 5.704604915246542, Index: 41, Point: core.NewPdfPoint(74.16588992664359, 81.00550169213861)}, {Distance: 7.175761206917839, Index: 58, Point: core.NewPdfPoint(70.09436264721477, 86.01779473250726)}, {Distance: 11.512743286703174, Index: 46, Point: core.NewPdfPoint(67.53711275495073, 90.39020510040851)}}},
		{query: core.NewPdfPoint(36.86160352788826, 13.756241445424322), expected: [3]NearestResult{{Distance: 8.660213312722579, Index: 13, Point: core.NewPdfPoint(35.803655214742406, 5.16089154047169)}, {Distance: 10.05626008778861, Index: 78, Point: core.NewPdfPoint(44.866927314279856, 7.669932380951572)}, {Distance: 11.012079230414543, Index: 88, Point: core.NewPdfPoint(43.03278016511083, 4.6357991376672185)}}},
		{query: core.NewPdfPoint(0.6887485300444696, 90.80035635580387), expected: [3]NearestResult{{Distance: 15.52506007098662, Index: 18, Point: core.NewPdfPoint(14.158071542303697, 98.52077265281578)}, {Distance: 20.774155245195235, Index: 68, Point: core.NewPdfPoint(10.747687004702266, 72.62390777764313)}, {Distance: 24.067856284005092, Index: 97, Point: core.NewPdfPoint(8.420107657979937, 68.0080815210619)}}},
		{query: core.NewPdfPoint(37.792743719353716, 82.90805938696141), expected: [3]NearestResult{{Distance: 4.253224348519968, Index: 66, Point: core.NewPdfPoint(34.43836768761709, 85.52303143204219)}, {Distance: 9.939820266481515, Index: 98, Point: core.NewPdfPoint(28.000497306010118, 81.20161949340638)}, {Distance: 13.848241320909453, Index: 10, Point: core.NewPdfPoint(24.41810676553939, 79.31739454102991)}}},
		{query: core.NewPdfPoint(28.425825069053435, 93.38421306346834), expected: [3]NearestResult{{Distance: 9.896910900999169, Index: 66, Point: core.NewPdfPoint(34.43836768761709, 85.52303143204219)}, {Distance: 12.190015988477223, Index: 98, Point: core.NewPdfPoint(28.000497306010118, 81.20161949340638)}, {Distance: 14.62659185673608, Index: 10, Point: core.NewPdfPoint(24.41810676553939, 79.31739454102991)}}},
		{query: core.NewPdfPoint(48.45367432788545, 3.2912155909772367), expected: [3]NearestResult{{Distance: 5.585158765696226, Index: 88, Point: core.NewPdfPoint(43.03278016511083, 4.6357991376672185)}, {Distance: 5.660204489805277, Index: 78, Point: core.NewPdfPoint(44.866927314279856, 7.669932380951572)}, {Distance: 5.982992148213537, Index: 43, Point: core.NewPdfPoint(47.945049322341525, 9.252548969452246)}}},
		{query: core.NewPdfPoint(41.681865022692435, 82.6132344499873), expected: [3]NearestResult{{Distance: 7.806098392924558, Index: 66, Point: core.NewPdfPoint(34.43836768761709, 85.52303143204219)}, {Distance: 10.876229338256989, Index: 86, Point: core.NewPdfPoint(52.44257355334966, 84.19422038812934)}, {Distance: 12.7306029372109, Index: 36, Point: core.NewPdfPoint(49.53804959253269, 92.63065004232082)}}},
		{query: core.NewPdfPoint(7.175621818812649, 4.404267177041543), expected: [3]NearestResult{{Distance: 4.700030158625489, Index: 17, Point: core.NewPdfPoint(8.254965696884042, 8.97868520394496)}, {Distance: 6.884800728051085, Index: 4, Point: core.NewPdfPoint(0.7742372570967326, 6.938583013317256)}, {Distance: 8.81749451695195, Index: 64, Point: core.NewPdfPoint(15.264247773398054, 0.8939481426395335)}}},
		{query: core.NewPdfPoint(67.06375058678753, 2.0706332172362285), expected: [3]NearestResult{{Distance: 3.807121420419626, Index: 76, Point: core.NewPdfPoint(70.36073331446983, 3.9743344551363857)}, {Distance: 4.597343034838145, Index: 92, Point: core.NewPdfPoint(69.38884149074713, 6.0366785105100185)}, {Distance: 7.521027118921692, Index: 48, Point: core.NewPdfPoint(71.56130305322326, 8.098723357173876)}}},
		{query: core.NewPdfPoint(21.17990951913106, 20.77232612239539), expected: [3]NearestResult{{Distance: 2.949689300873121, Index: 50, Point: core.NewPdfPoint(18.3788032374675, 19.84795034896214)}, {Distance: 2.9500005402388694, Index: 99, Point: core.NewPdfPoint(24.103817108644833, 21.163820177769633)}, {Distance: 3.683726578350228, Index: 38, Point: core.NewPdfPoint(24.636942156010523, 22.044633809864933)}}},
		{query: core.NewPdfPoint(83.4783024984307, 24.393830521827052), expected: [3]NearestResult{{Distance: 3.0961634442512573, Index: 8, Point: core.NewPdfPoint(80.59123120177551, 25.51233282470693)}, {Distance: 3.7225758420518424, Index: 15, Point: core.NewPdfPoint(80.11998397458856, 25.99984035749021)}, {Distance: 3.768005736156072, Index: 94, Point: core.NewPdfPoint(80.69475374791655, 21.854199922464968)}}},
		{query: core.NewPdfPoint(67.94892501031262, 58.40285832426694), expected: [3]NearestResult{{Distance: 7.002211460552816, Index: 59, Point: core.NewPdfPoint(71.14819197800337, 52.17424240997747)}, {Distance: 10.457617375062723, Index: 23, Point: core.NewPdfPoint(59.31519405414486, 64.30374393065894)}, {Distance: 12.58403246559697, Index: 30, Point: core.NewPdfPoint(79.77968787046065, 62.69132229661655)}}},
		{query: core.NewPdfPoint(70.30054809098338, 38.229402533712175), expected: [3]NearestResult{{Distance: 3.3827550961350985, Index: 87, Point: core.NewPdfPoint(67.0752294529232, 39.24937886844233)}, {Distance: 6.676593719196517, Index: 63, Point: core.NewPdfPoint(72.74747856539602, 44.441442049411116)}, {Distance: 8.77548579112744, Index: 74, Point: core.NewPdfPoint(75.03458936225277, 45.61844875850991)}}},
		{query: core.NewPdfPoint(6.666591704963154, 10.002159844611135), expected: [3]NearestResult{{Distance: 1.8895587522745774, Index: 17, Point: core.NewPdfPoint(8.254965696884042, 8.97868520394496)}, {Distance: 5.383387494014664, Index: 26, Point: core.NewPdfPoint(11.228694290441743, 12.860145308234605)}, {Distance: 6.641185431873729, Index: 4, Point: core.NewPdfPoint(0.7742372570967326, 6.938583013317256)}}},
		{query: core.NewPdfPoint(37.737895538821185, 56.99302071734109), expected: [3]NearestResult{{Distance: 5.280595736713236, Index: 60, Point: core.NewPdfPoint(43.017566199578106, 57.09185921458071)}, {Distance: 6.8049464384399085, Index: 35, Point: core.NewPdfPoint(31.112041569584935, 55.44209995586027)}, {Distance: 7.120979010551258, Index: 84, Point: core.NewPdfPoint(42.195134616764996, 62.54652068319705)}}},
		{query: core.NewPdfPoint(46.09269243297545, 55.527895517428725), expected: [3]NearestResult{{Distance: 3.449983158993363, Index: 60, Point: core.NewPdfPoint(43.017566199578106, 57.09185921458071)}, {Distance: 8.028203793393617, Index: 84, Point: core.NewPdfPoint(42.195134616764996, 62.54652068319705)}, {Distance: 14.481241345005133, Index: 71, Point: core.NewPdfPoint(52.97863363490906, 68.26721122337617)}}},
		{query: core.NewPdfPoint(53.61166019135128, 29.402010229671795), expected: [3]NearestResult{{Distance: 7.933994891088919, Index: 96, Point: core.NewPdfPoint(46.898436066324166, 25.173429288152636)}, {Distance: 10.207161689017788, Index: 85, Point: core.NewPdfPoint(49.61310425442538, 20.010647167526496)}, {Distance: 11.287466296847745, Index: 89, Point: core.NewPdfPoint(49.42598643546282, 18.91930873569896)}}},
		{query: core.NewPdfPoint(44.192808857614644, 34.74442504823799), expected: [3]NearestResult{{Distance: 6.465911940748226, Index: 49, Point: core.NewPdfPoint(37.970724110053844, 36.503310666941424)}, {Distance: 9.469778598317067, Index: 53, Point: core.NewPdfPoint(37.57869607184197, 27.967241358959505)}, {Distance: 9.946073518332808, Index: 96, Point: core.NewPdfPoint(46.898436066324166, 25.173429288152636)}}},
		{query: core.NewPdfPoint(87.18609487820613, 36.034544097040076), expected: [3]NearestResult{{Distance: 4.658256014904542, Index: 79, Point: core.NewPdfPoint(90.03637785608008, 39.71900281528301)}, {Distance: 4.96949613369587, Index: 72, Point: core.NewPdfPoint(84.97188896260866, 40.483497230772926)}, {Distance: 5.333717325854641, Index: 33, Point: core.NewPdfPoint(84.80993603153726, 40.80973096746926)}}},
		{query: core.NewPdfPoint(72.26195709778688, 9.598162829344847), expected: [3]NearestResult{{Distance: 1.5467587665908336, Index: 19, Point: core.NewPdfPoint(70.71899728005727, 9.70650341683359)}, {Distance: 1.6550633887763093, Index: 48, Point: core.NewPdfPoint(71.56130305322326, 8.098723357173876)}, {Distance: 4.575911258396712, Index: 92, Point: core.NewPdfPoint(69.38884149074713, 6.0366785105100185)}}},
		{query: core.NewPdfPoint(34.1612429164763, 19.22222038671292), expected: [3]NearestResult{{Distance: 8.368264652666655, Index: 2, Point: core.NewPdfPoint(35.583597244228336, 27.46872033901967)}, {Distance: 9.389056282404074, Index: 53, Point: core.NewPdfPoint(37.57869607184197, 27.967241358959505)}, {Distance: 9.933696316427113, Index: 38, Point: core.NewPdfPoint(24.636942156010523, 22.044633809864933)}}},
		{query: core.NewPdfPoint(71.38391153141247, 13.520311281559493), expected: [3]NearestResult{{Distance: 3.871335866429107, Index: 19, Point: core.NewPdfPoint(70.71899728005727, 9.70650341683359)}, {Distance: 5.424489227001405, Index: 48, Point: core.NewPdfPoint(71.56130305322326, 8.098723357173876)}, {Distance: 7.745002512529356, Index: 92, Point: core.NewPdfPoint(69.38884149074713, 6.0366785105100185)}}},
		{query: core.NewPdfPoint(98.937208620892, 67.26912373356589), expected: [3]NearestResult{{Distance: 8.43618896742764, Index: 37, Point: core.NewPdfPoint(92.45607820110665, 72.66951564303699)}, {Distance: 14.035035311182662, Index: 55, Point: core.NewPdfPoint(84.9666796464414, 68.61319827025679)}, {Distance: 14.3852303261261, Index: 39, Point: core.NewPdfPoint(99.76982890678676, 52.90800974490838)}}},
		{query: core.NewPdfPoint(89.80071239640789, 2.42574224851011), expected: [3]NearestResult{{Distance: 5.724587358789773, Index: 28, Point: core.NewPdfPoint(90.87481745321885, 8.048659530077462)}, {Distance: 8.568651429497155, Index: 54, Point: core.NewPdfPoint(97.7460348272766, 5.634112362720456)}, {Distance: 13.056548614184756, Index: 11, Point: core.NewPdfPoint(98.12630738875605, 12.483472098628233)}}},
		{query: core.NewPdfPoint(4.9811550769006345, 11.349911276229651), expected: [3]NearestResult{{Distance: 4.042344500583644, Index: 17, Point: core.NewPdfPoint(8.254965696884042, 8.97868520394496)}, {Distance: 6.095734130172781, Index: 4, Point: core.NewPdfPoint(0.7742372570967326, 6.938583013317256)}, {Distance: 6.4274841933807805, Index: 26, Point: core.NewPdfPoint(11.228694290441743, 12.860145308234605)}}},
	}
	for _, tc := range tests {
		results := tree.FindNearestNeighbours(tc.query, 3, Euclidean)
		if len(results) != 3 { t.Errorf("expected 3 results, got %d", len(results)); continue }
		for i := 0; i < 3; i++ {
			r := results[i]; e := tc.expected[i]
			assertNearKdTree(t, r.Distance, e.Distance, 0.000001)
			if r.Index != e.Index { t.Errorf("[%d] idx: got %d, want %d", i, r.Index, e.Index) }
			if math.Abs(r.Point.X-e.Point.X)>0.001||math.Abs(r.Point.Y-e.Point.Y)>0.001 {
				t.Errorf("[%d] pt: got %v, want %v", i, r.Point, e.Point)
			}
		}
	}
}
