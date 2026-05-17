//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	dla "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/word_extractor"
)

func TestGetBlocks(t *testing.T) {
	name := "Random 2 Columns Lists Hyph - Justified.pdf"
	expected := []string{
		"Random Big Title",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. In sodales gravida felis, in rhoncus velit rutrum at. Curabitur hendrerit dapibus nulla, ut hendrerit diam imperdiet quis. Pellentesque id neque ali- quam, pulvinar neque in, vulputate elit. Pel- lentesque ut erat sit amet massa suscipit ullamcor- per. Sed porttitor viverra convallis. Duis vitae sem- per metus. Pellentesque eros purus, egestas eget velit eget, elementum aliquet velit. Suspendisse potenti. Nulla vitae massa rutrum, blandit erat vi- tae, aliquet arcu.",
		"Aenean feugiat leo sed enim sodales vehicula. Sus- pendisse tempus hendrerit magna sagittis dictum. Duis ultrices dapibus egestas. Cras eu felis eu lectus suscipit pharetra at at lacus. Nulla facilisi. Proin in- terdum faucibus elit nec rhoncus. Proin sodaless metus sed tincidunt hendrerit.",
		"Donec ultricies cursus odio sed rutrum. Nam ven- enatis metus vitae elementum scelerisque. Ali- quam tempor sapien at turpis posuere eleifend. Sed placerat posuere nunc vel efficitur. Quisque auctor felis vel lectus dictum fringilla. Quisque vo- lutpat pulvinar© elit. Aliquam ultrices feugiat ali- quam. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sus- pendisse imperdiet ex lorem, porta bibendum pu- rus ultricies id.",
		"Integer vel lacus sapien. Nam sodales ante eu risus facilisis placerat. Aliquam suscipit pulvinar ultricies. Aenean pulvinar, ex ac fermentum egestas, erat nisi feugiat velit, vitae suscipit tellus odio vitae quam. Morbi elementum sem in elit posuere, non",
		"• Duis leo enim, convallis sit amet orci eget, condimentum mattis mi ; • Etiam dolor erat, maximus nec mi sed, con- vallis convallis orci ; • Morbi viverra diam in diam cursus, vitae aliquet velit tempus ; • Donec at nisi fermentum, ultricies odio eget, egestas massa at nisi fermentum, ul- tricies odio eget, egestas massa.",
		"Lorem Ipsum text with lists",
		"rhoncus magna fringilla. Phasellus cursus in dolor laoreet rutrum. Curabitur tincidunt risus ullamcor- per, vehicula velit at, pulvinar metus.",
		"Donec quis ante leo. Vivamus pharetra, nisl ac vehi- cula tempor, tellus lacus aliquam sapien, eu congue nibh quam sit amet odio. Quisque metus arcu, sem- per nec consequat eu, pellentesque vel sem. Sed purus risus, tincidunt¹ sit amet dictum vitae, euis- mod id nibh. Praesent ultrices libero quis enim porta, sit amet pellentesque augue pretium. Viva- mus nec molestie nunc. Donec finibus enim nec tel- lus laoreet elementum. Curabitur efficitur placerat dolor et semper.",
		"Morbi laoreet dui eu tortor luctus, nec ultrices do- lor ullamcorper. Ut gravida sed nisl a efficitur. In tincidunt orci a condimentum semper. Suspendisse scelerisque fermentum lacinia. Vestibulum sit amet ornare tellus, aliquet euismod mauris. Cras suscipit venenatis ultrices. Sed diam erat, aliquet a tellus ut, viverra 12º ongue magna. Cras id justo tortor. Mauris in tortor vulputate, pellentesque nisl ac, facilisis ligula. Class aptent taciti² sociosqu ad li- tora torquent per conubia nostra³, per inceptos himenaeos. Aliquam eget dolor turpis. Mauris id molestie tellus. Sed elementum molestie nisi, at ali- quet sem vehicula nec. Morbi tempus nulla enim, a vulputate magna €51 luctus £66 eu. Fusce sodales, libero quis suscipit ultrices, metus erat auctor urna, sit amet dictum arcu tortor eu metus.",
		"Morbi vestibulum varius ipsum nec molestie. Proin auctor efficitur diam ut luctus. Phasellus cursus maximus ultricies. Mauris eu neque ut sem semper tempus. Curabitur non lorem eu nunc lobortis vi- verra at in diam. Pellentesque euismod purus a leo lobortis tempor. Maecenas mollis ligula at sem sus- cipit fringilla. Mauris sollicitudin tincidunt lectus id tempor. Etiam ut nisi est.",
		"1. Ut volutpat, velit at interdum consectetur, nisl lorem consequat mauris, feugiat dignissim tellus massa ut nisl. 2. Praesent at est nisi. Pellentesque rutrum lorem sed dui accumsan gravida. 3. Pellentesque dictum nisl vitae urna luctus, congue pulvinar mi congue.",
	}

	docPath := filepath.Join(dlaDocRoot, name)
	if abs, err := filepath.Abs(docPath); err == nil {
		docPath = abs
	}
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		t.Skipf("Test document not found: %s", docPath)
	}

	parsingOpts := &content.ParsingOptions{UseLenientParsing: true}
	doc, err := pdfpig.OpenFile(docPath, parsingOpts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", docPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	extractor := word_extractor.NewNearestNeighbourWordExtractor()
	words := extractor.GetWords(page.Letters())

	opts := page_segmenter.DefaultRecursiveXYCutOptions()
	opts.SetMinimumWidth(page.Width() / 3.0)
	opts.SetLineSeparator(" ")

	segmenter := page_segmenter.NewRecursiveXYCut(opts)
	blocks := segmenter.GetBlocks(words)

	if len(blocks) != len(expected) {
		t.Fatalf("expected %d blocks, got %d", len(expected), len(blocks))
	}

	ordered := make([]*dla.TextBlock, len(blocks))
	copy(ordered, blocks)
	sort.SliceStable(ordered, func(i, j int) bool {
		bi, bj := ordered[i].BoundingBox.BottomLeft, ordered[j].BoundingBox.BottomLeft
		if bi.X != bj.X {
			return bi.X < bj.X
		}
		return bi.Y > bj.Y
	})

	for i, expText := range expected {
		actual := strings.TrimSpace(ordered[i].Text)
		expTrimmed := strings.TrimSpace(expText)
		if actual != expTrimmed {
			t.Errorf("block %d:\nexpected: %q\n  actual: %q", i, expTrimmed, actual)
		}
	}
}

func TestGetWords(t *testing.T) {
	name := "Random 2 Columns Lists Hyph - Justified.pdf"

	docPath := filepath.Join(dlaDocRoot, name)
	if abs, err := filepath.Abs(docPath); err == nil {
		docPath = abs
	}
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		t.Skipf("Test document not found: %s", docPath)
	}

	parsingOpts := &content.ParsingOptions{UseLenientParsing: true}
	doc, err := pdfpig.OpenFile(docPath, parsingOpts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", docPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	extractor := word_extractor.NewNearestNeighbourWordExtractor()
	words := extractor.GetWords(page.Letters())

	if len(words) == 0 {
		t.Fatal("expected at least one word, got 0")
	}

	for i, word := range words {
		if word.Text == "" {
			t.Errorf("word %d has empty text", i)
		}
		if len(word.Letters) == 0 {
			t.Errorf("word %d (%q) has no letters", i, word.Text)
		}
	}
}
