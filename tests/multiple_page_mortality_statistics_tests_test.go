//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

const mortalityStatisticsDocName = "Multiple Page - from Mortality Statistics.pdf"

func resolveMortalityStatisticsPath() string {
	return filepath.Join(integrationDocRoot, mortalityStatisticsDocName)
}

// TestMortalityStatisticsHasCorrectNumberOfPages verifies the document has 6 pages.
func TestMortalityStatisticsHasCorrectNumberOfPages(t *testing.T) {
	path := resolveMortalityStatisticsPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.NumberOfPages(); got != 6 {
		t.Errorf("NumberOfPages = %d, want 6", got)
	}
}

// TestMortalityStatisticsHasCorrectVersion verifies the PDF version is 1.7.
func TestMortalityStatisticsHasCorrectVersion(t *testing.T) {
	path := resolveMortalityStatisticsPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if got := doc.Version(); got != 1.7 {
		t.Errorf("Version = %.1f, want 1.7", got)
	}
}

// TestMortalityStatisticsGetsFirstPageContent verifies text content and page size on page 1.
func TestMortalityStatisticsGetsFirstPageContent(t *testing.T) {
	path := resolveMortalityStatisticsPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("page is not *content.Page")
	}

	text := page.Text()

	expectedFragments := []string{
		"Mortality Statistics: Metadata",
		"Notification to the registrar by the coroner that he does not consider it necessary to hold an inquest – no post-mortem held (Form 100A – salmon pink)",
		"Presumption of death certificate",
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(text, fragment) {
			t.Errorf("Page 1 text does not contain %q", fragment)
		}
	}

	if size := page.Size(); size != content.PageSizeLetter {
		t.Errorf("Page 1 Size = %v, want PageSizeLetter", size)
	}
}

// TestMortalityStatisticsGetsPagesContent verifies text on page 6.
func TestMortalityStatisticsGetsPagesContent(t *testing.T) {
	path := resolveMortalityStatisticsPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(6)
	if err != nil {
		t.Fatalf("GetPage(6): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("page is not *content.Page")
	}

	expected := "Up to 1992, publications gave numbers of deaths registered in the period concerned. From 1993 to 2005, the figures in annual reference volumes relate to the number of deaths that occurred in the reference period. From 2006 onwards, all tables in Series DR are based on deaths registered in a calendar period. More details on these changes can be found in the publication Mortality Statistics: Deaths Registered in 2006 (ONS, 2008)"

	if !strings.Contains(page.Text(), expected) {
		t.Errorf("Page 6 text does not contain expected fragment about death statistics methodology")
	}
}

// TestMortalityStatisticsLettersHaveCorrectPositionsPdfBox verifies letter positions on page 1.
func TestMortalityStatisticsLettersHaveCorrectPositionsPdfBox(t *testing.T) {
	path := resolveMortalityStatisticsPath()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("page is not *content.Page")
	}

	letters := page.Letters()
	positions := getPdfBoxPositions()

	for i := 0; i < len(positions); i++ {
		if i >= len(letters) {
			t.Fatalf("Expected at least %d letters, got %d", len(positions), len(letters))
		}

		positions[i].AssertWithinTolerance(t, letters[i], false)
	}
}

func getPdfBoxPositions() []*testutil.AssertablePositionData {
	const data = "390.6\t741.12\t8.29668\tM\t9.96\tArialMT\t14.601361\n" +
		"398.87976\t741.12\t5.5377607\to\t9.96\tArialMT\t11.0556\n" +
		"404.39957\t741.12\t3.31668\tr\t9.96\tArialMT\t10.816561\n" +
		"407.76007\t741.12\t2.7688804\tt\t9.96\tArialMT\t14.41212\n" +
		"410.52\t741.12\t5.5377607\ta\t9.96\tArialMT\t11.0556\n" +
		"416.1603\t741.12\t2.2111201\tl\t9.96\tArialMT\t14.601361\n" +
		"418.32162\t741.12\t2.2111201\ti\t9.96\tArialMT\t14.601361\n" +
		"420.48294\t741.12\t2.7688804\tt\t9.96\tArialMT\t14.41212\n" +
		"423.4839\t741.12\t4.98\ty\t9.96\tArialMT\t14.87028\n" +
		"428.28\t741.12\t2.7688804\t \t9.96\tArialMT\t0.0\n" +
		"431.16\t741.12\t6.6433206\tS\t9.96\tArialMT\t15.09936\n" +
		"437.75952\t741.12\t2.7688804\tt\t9.96\tArialMT\t14.41212\n" +
		"440.52045\t741.12\t5.5377607\ta\t9.96\tArialMT\t11.0556\n" +
		"446.0383\t741.12\t2.7688804\tt\t9.96\tArialMT\t14.41212\n" +
		"448.91772\t741.12\t2.2111201\ti\t9.96\tArialMT\t14.601361\n" +
		"451.07803\t741.12\t4.98\ts\t9.96\tArialMT\t11.0556\n" +
		"456.1178\t741.12\t2.7688804\tt\t9.96\tArialMT\t14.41212\n" +
		"458.87772\t741.12\t2.2111201\ti\t9.96\tArialMT\t14.601361\n" +
		"461.03802\t741.12\t4.98\tc\t9.96\tArialMT\t11.0556\n" +
		"466.0778\t741.12\t4.98\ts\t9.96\tArialMT\t11.0556\n" +
		"471.11755\t741.12\t2.7688804\t:\t9.96\tArialMT\t10.57752\n" +
		"473.87747\t741.12\t2.7688804\t \t9.96\tArialMT\t0.0\n" +
		"476.63638\t741.12\t8.29668\tM\t9.96\tArialMT\t14.601361\n" +
		"485.03564\t741.12\t5.5377607\te\t9.96\tArialMT\t11.0556\n" +
		"490.5535\t741.12\t2.7688804\tt\t9.96\tArialMT\t14.41212\n" +
		"493.31342\t741.12\t5.5377607\ta\t9.96\tArialMT\t11.0556\n" +
		"498.83127\t741.12\t5.5377607\td\t9.96\tArialMT\t14.840402\n" +
		"504.47064\t741.12\t5.5377607\ta\t9.96\tArialMT\t11.0556\n" +
		"509.9885\t741.12\t2.7688804\tt\t9.96\tArialMT\t14.41212\n" +
		"512.8689\t741.12\t5.5377607\ta\t9.96\tArialMT\t11.0556\n" +
		"518.4\t741.12\t2.7688804\t \t9.96\tArialMT\t0.0\n" +
		"571.92\t741.36\t2.5020003\t \t9.0\tArialMT\t0.0\n" +
		"90.0\t721.44\t3.0\t \t12.0\tTimesNewRomanPSMT\t0.0\n" +
		"90.0\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"126.001434\t708.24\t6.1382403\t4\t11.04\tArialMT\t16.18464\n" +
		"132.1176\t708.24\t3.0691202\t.\t11.04\tArialMT\t2.2632\n" +
		"135.23862\t708.24\t6.1382403\t3\t11.04\tArialMT\t16.53792\n" +
		"141.23663\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"161.99184\t708.24\t7.3636804\tE\t11.04\tArialMT\t16.18464\n" +
		"169.31136\t708.24\t6.1382403\tn\t11.04\tArialMT\t11.989441\n" +
		"175.42752\t708.24\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"178.54742\t708.24\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"182.2679\t708.24\t5.52\ty\t11.04\tArialMT\t16.48272\n" +
		"187.66757\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"190.79189\t708.24\t6.1382403\to\t11.04\tArialMT\t12.2544\n" +
		"196.79213\t708.24\t3.0691202\tf\t11.04\tArialMT\t16.46064\n" +
		"199.91203\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"203.03635\t708.24\t6.1382403\td\t11.04\tArialMT\t16.449602\n" +
		"209.15251\t708.24\t6.1382403\ta\t11.04\tArialMT\t12.2544\n" +
		"215.26868\t708.24\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"218.38858\t708.24\t6.1382403\ta\t11.04\tArialMT\t12.2544\n" +
		"224.38773\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"227.39279\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"233.9947\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"269.99615\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"305.9976\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"341.99902\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"378.00046\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"414.0019\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"450.00333\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"486.00476\t708.24\t6.1382403\t3\t11.04\tArialMT\t16.53792\n" +
		"492.1209\t708.24\t6.1382403\t6\t11.04\tArialMT\t16.52688\n" +
		"498.23706\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"521.9951\t708.24\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"162.0\t695.52\t3.6984\t\u2022\t8.04\tLYVNDC+SymbolMT\t0.008040001\n" +
		"165.72252\t695.52\t2.2351203\t \t8.04\tArialMT\t0.0\n" +
		"180.0\t695.52\t7.9708805\tR\t11.04\tArialMT\t16.18464\n" +
		"187.9201\t695.52\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"194.03627\t695.52\t6.1382403\tg\t11.04\tArialMT\t16.74768\n" +
		"200.27719\t695.52\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"202.6762\t695.52\t5.52\ts\t11.04\tArialMT\t12.2544\n" +
		"208.1962\t695.52\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"211.19687\t695.52\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"214.91624\t695.52\t6.1382403\ta\t11.04\tArialMT\t12.2544\n" +
		"221.03241\t695.52\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"224.15231\t695.52\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"226.55241\t695.52\t6.1382403\to\t11.04\tArialMT\t12.2544\n" +
		"232.66858\t695.52\t6.1382403\tn\t11.04\tArialMT\t11.989441\n" +
		"238.78474\t695.52\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"241.90906\t695.52\t7.3636804\tS\t11.04\tArialMT\t16.73664\n" +
		"249.22858\t695.52\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"255.22882\t695.52\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"258.9493\t695.52\t5.52\tv\t11.04\tArialMT\t11.724481\n" +
		"264.34897\t695.52\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"266.74905\t695.52\t5.52\tc\t11.04\tArialMT\t12.2544\n" +
		"272.26907\t695.52\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"278.38522\t695.52\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"281.50955\t695.52\t7.3636804\tS\t11.04\tArialMT\t16.73664\n" +
		"288.8302\t695.52\t6.1382403\to\t11.04\tArialMT\t12.2544\n" +
		"294.8293\t695.52\t3.0691202\tf\t11.04\tArialMT\t16.46064\n" +
		"298.06955\t695.52\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"301.0691\t695.52\t7.9708805\tw\t11.04\tArialMT\t11.724481\n" +
		"308.86993\t695.52\t6.1382403\ta\t11.04\tArialMT\t12.2544\n" +
		"314.98608\t695.52\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"318.70657\t695.52\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"324.82272\t695.52\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"327.94705\t695.52\t3.67632\t(\t11.04\tArialMT\t21.21888\n" +
		"331.66754\t695.52\t7.9708805\tR\t11.04\tArialMT\t16.18464\n" +
		"339.58762\t695.52\t7.3636804\tS\t11.04\tArialMT\t16.73664\n" +
		"346.90714\t695.52\t7.3636804\tS\t11.04\tArialMT\t16.73664\n" +
		"354.22778\t695.52\t3.67632\t)\t11.04\tArialMT\t21.21888\n" +
		"357.9648\t695.52\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"162.0\t682.92\t3.6984\t\u2022\t8.04\tLYVNDC+SymbolMT\t0.008040001\n" +
		"165.72252\t682.92\t2.2351203\t \t8.04\tArialMT\t0.0\n" +
		"180.0\t682.92\t7.9708805\tR\t11.04\tArialMT\t16.18464\n" +
		"187.9201\t682.92\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"194.03627\t682.92\t6.1382403\tg\t11.04\tArialMT\t16.74768\n" +
		"200.27719\t682.92\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"202.6762\t682.92\t5.52\ts\t11.04\tArialMT\t12.2544\n" +
		"208.1962\t682.92\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"211.19687\t682.92\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"214.91624\t682.92\t6.1382403\ta\t11.04\tArialMT\t12.2544\n" +
		"221.03241\t682.92\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"224.15231\t682.92\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"226.55241\t682.92\t6.1382403\to\t11.04\tArialMT\t12.2544\n" +
		"232.66858\t682.92\t6.1382403\tn\t11.04\tArialMT\t11.989441\n" +
		"238.78474\t682.92\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"241.78763\t682.92\t8.589121\tO\t11.04\tArialMT\t16.74768\n" +
		"250.42754\t682.92\t6.1382403\tn\t11.04\tArialMT\t11.989441\n" +
		"256.5437\t682.92\t2.4508803\tl\t11.04\tArialMT\t16.18464\n" +
		"258.9427\t682.92\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"261.34277\t682.92\t6.1382403\tn\t11.04\tArialMT\t11.989441\n" +
		"267.45892\t682.92\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"273.57507\t682.92\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"276.69498\t682.92\t3.67632\t(\t11.04\tArialMT\t21.21888\n" +
		"280.41547\t682.92\t7.9708805\tR\t11.04\tArialMT\t16.18464\n" +
		"288.2152\t682.92\t8.589121\tO\t11.04\tArialMT\t16.74768\n" +
		"296.7359\t682.92\t7.9708805\tN\t11.04\tArialMT\t16.18464\n" +
		"304.65598\t682.92\t3.67632\t)\t11.04\tArialMT\t21.21888\n" +
		"308.3952\t682.92\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"89.990875\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"125.99231\t670.2019\t6.1382403\t4\t11.04\tArialMT\t16.18464\n" +
		"132.10847\t670.2019\t3.0691202\t.\t11.04\tArialMT\t2.2632\n" +
		"135.22838\t670.2019\t6.1382403\t4\t11.04\tArialMT\t16.18464\n" +
		"141.22751\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"161.98271\t670.2019\t8.589121\tO\t11.04\tArialMT\t16.74768\n" +
		"170.62262\t670.2019\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"173.74252\t670.2019\t6.1382403\th\t11.04\tArialMT\t16.18464\n" +
		"179.85869\t670.2019\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"185.85783\t670.2019\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"189.57831\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"192.70262\t670.2019\t5.52\tc\t11.04\tArialMT\t12.2544\n" +
		"198.22263\t670.2019\t6.1382403\th\t11.04\tArialMT\t16.18464\n" +
		"204.33879\t670.2019\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"210.33904\t670.2019\t5.52\tc\t11.04\tArialMT\t12.2544\n" +
		"215.7398\t670.2019\t5.52\tk\t11.04\tArialMT\t16.18464\n" +
		"221.37904\t670.2019\t5.52\ts\t11.04\tArialMT\t12.2544\n" +
		"226.89905\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"229.8986\t670.2019\t9.196321\tm\t11.04\tArialMT\t11.989441\n" +
		"239.13908\t670.2019\t6.1382403\ta\t11.04\tArialMT\t12.2544\n" +
		"245.25525\t670.2019\t6.1382403\td\t11.04\tArialMT\t16.449602\n" +
		"251.37141\t670.2019\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"257.37167\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"260.496\t670.2019\t6.1382403\tb\t11.04\tArialMT\t16.449602\n" +
		"266.61215\t670.2019\t5.52\ty\t11.04\tArialMT\t16.48272\n" +
		"272.01294\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"275.13727\t670.2019\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"278.25717\t670.2019\t6.1382403\th\t11.04\tArialMT\t16.18464\n" +
		"284.2563\t670.2019\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"290.37244\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"293.49677\t670.2019\t7.9708805\tR\t11.04\tArialMT\t16.18464\n" +
		"301.41684\t670.2019\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"307.533\t670.2019\t6.1382403\tg\t11.04\tArialMT\t16.74768\n" +
		"313.7728\t670.2019\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"316.17288\t670.2019\t5.52\ts\t11.04\tArialMT\t12.2544\n" +
		"321.57367\t670.2019\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"324.69357\t670.2019\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"328.41293\t670.2019\t6.1382403\ta\t11.04\tArialMT\t12.2544\n" +
		"334.41318\t670.2019\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"337.53308\t670.2019\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"339.93317\t670.2019\t6.1382403\to\t11.04\tArialMT\t12.2544\n" +
		"346.04932\t670.2019\t6.1382403\tn\t11.04\tArialMT\t11.989441\n" +
		"352.16547\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"355.2898\t670.2019\t7.3636804\tS\t11.04\tArialMT\t16.73664\n" +
		"362.6093\t670.2019\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"368.72546\t670.2019\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"372.44482\t670.2019\t5.52\tv\t11.04\tArialMT\t11.724481\n" +
		"377.8456\t670.2019\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"380.2457\t670.2019\t5.52\tc\t11.04\tArialMT\t12.2544\n" +
		"385.76572\t670.2019\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"391.88187\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"394.90463\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"413.98175\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"449.9832\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"485.98462\t670.2019\t6.1382403\t3\t11.04\tArialMT\t16.53792\n" +
		"492.10077\t670.2019\t6.1382403\t7\t11.04\tArialMT\t15.97488\n" +
		"498.21692\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"501.2198\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"521.975\t670.2019\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"162.0\t657.6\t3.6984\t\u2022\t8.04\tLYVNDC+SymbolMT\t0.008040001\n" +
		"165.72252\t657.6\t2.2351203\t \t8.04\tArialMT\t0.0\n" +
		"180.0\t657.6\t7.3636804\tA\t11.04\tArialMT\t16.18464\n" +
		"187.31952\t657.6\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"190.44383\t657.6\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"193.56374\t657.6\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"196.68805\t657.6\t6.1382403\th\t11.04\tArialMT\t16.18464\n" +
		"202.80862\t657.6\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"208.80885\t657.6\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"211.81174\t657.6\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"214.93605\t657.6\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"217.33614\t657.6\t9.196321\tm\t11.04\tArialMT\t11.989441\n" +
		"226.57552\t657.6\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"232.69609\t657.6\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"235.816\t657.6\t6.1382403\to\t11.04\tArialMT\t12.2544\n" +
		"241.81622\t657.6\t3.0691202\tf\t11.04\tArialMT\t16.46064\n" +
		"244.94054\t657.6\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"247.94342\t657.6\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"251.6639\t657.6\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"257.66302\t657.6\t6.1382403\tg\t11.04\tArialMT\t16.74768\n" +
		"263.90283\t657.6\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"266.30292\t657.6\t5.52\ts\t11.04\tArialMT\t12.2544\n" +
		"271.8229\t657.6\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"274.82358\t657.6\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"278.54294\t657.6\t6.1382403\ta\t11.04\tArialMT\t12.2544\n" +
		"284.6635\t657.6\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"287.78784\t657.6\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"290.18683\t657.6\t6.1382403\to\t11.04\tArialMT\t12.2544\n" +
		"296.3074\t657.6\t6.1382403\tn\t11.04\tArialMT\t11.989441\n" +
		"302.27905\t657.6\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"162.0\t645.0\t3.6984\t\u2022\t8.04\tLYVNDC+SymbolMT\t0.008040001\n" +
		"165.72252\t645.0\t2.2351203\t \t8.04\tArialMT\t0.0\n" +
		"180.0\t645.0\t7.3636804\tB\t11.04\tArialMT\t16.18464\n" +
		"187.31952\t645.0\t5.52\ty\t11.04\tArialMT\t16.48272\n" +
		"192.72029\t645.0\t3.0691202\t \t11.04\tArialMT\t0.0\n" +
		"195.8446\t645.0\t5.52\ts\t11.04\tArialMT\t12.2544\n" +
		"201.36461\t645.0\t6.1382403\tu\t11.04\tArialMT\t11.989441\n" +
		"207.48077\t645.0\t6.1382403\tp\t11.04\tArialMT\t16.48272\n" +
		"213.59694\t645.0\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"219.7131\t645.0\t3.67632\tr\t11.04\tArialMT\t11.989441\n" +
		"223.43358\t645.0\t2.4508803\ti\t11.04\tArialMT\t16.18464\n" +
		"225.83368\t645.0\t6.1382403\tn\t11.04\tArialMT\t11.989441\n" +
		"231.94984\t645.0\t3.0691202\tt\t11.04\tArialMT\t15.97488\n" +
		"235.06975\t645.0\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"241.18591\t645.0\t6.1382403\tn\t11.04\tArialMT\t11.989441\n" +
		"247.30208\t645.0\t6.1382403\td\t11.04\tArialMT\t16.449602\n" +
		"253.41824\t645.0\t6.1382403\te\t11.04\tArialMT\t12.2544\n" +
		"259.5344\t645.0\t6.1382403\tn\t11.04\tArialMT\t11.989441\n" +
		"265.65054\t645.0\t3.0691202\tt\t11.04\tArialMT\t15.97488"

	var result []*testutil.AssertablePositionData
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		pd, err := testutil.ParseAssertablePositionData(line)
		if err != nil {
			panic(err)
		}
		result = append(result, pd)
	}

	return result
}
