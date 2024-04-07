package steam

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{"heading1", `# foo`, "[h1]foo[/h1]\n"},
		{"heading2", `## foo`, "[h2]foo[/h2]\n"},
		{"heading3", `### foo`, "[h3]foo[/h3]\n"},
		{"heading4", `#### foo`, "[h3]foo[/h3]\n"},
		{"heading5", `##### foo`, "[h3]foo[/h3]\n"},
		{"heading6", `###### foo`, "[h3]foo[/h3]\n"},
		{"quote", "> quoted", "[quote]\nquoted\n\n[/quote]\n"},
		{"strong", "**strong**", "[b]strong[/b]\n\n"},
		{"ordered list", `
1. one
2. two
3. three`, `[olist]
[*]one
[*]two
[*]three
[/olist]
`},
		{"unordered list", `
- one
- two
- three`, `[list]
[*]one
[*]two
[*]three
[/list]
`},
		{"table", `
|head1 | head2|
|------|------|
|cell01|cell02|
|cell11|cell12|
`, `[table]
[tr]
[th]head1[/th][th]head2[/th][/tr]
[tr]
[td]cell01[/td][td]cell02[/td][/tr]
[tr]
[td]cell11[/td][td]cell12[/td][/tr]
[/table]
`},
		{"nested", `***strongem***`, "[b][i]strongem[/i][/b]\n\n"},
		{"nested", `# *em* **strong** ~~strike~~`, "[h1][i]em[/i] [b]strong[/b] [strike]strike[/strike][/h1]\n"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("name:%s", tc.name), func(t *testing.T) {
			t.Parallel()
			res := Run([]byte(tc.input))
			if string(res) != tc.want {
				t.Errorf("error producing output \n `%s` != `%s`", res, tc.want)
			}
		})
	}
}

func TestRunLarge(t *testing.T) {
	input := []byte(`
# heading one

## heading two

### heading three

#### heading four

##### heading five

> quote!
> very **strong**, there shall not be ~~*emphasis*~~

this is a simple paragraph with *strong* **arguments**!.

---

Link: [foo](https://example.com)
Image: ![](https://fastly.picsum.photos/id/133/200/300.jpg)
`)
	expected := `[h1]heading one[/h1]
[h2]heading two[/h2]
[h3]heading three[/h3]
[h3]heading four[/h3]
[h3]heading five[/h3]
[quote]
quote!
very [b]strong[/b], there shall not be [strike][i]emphasis[/i][/strike]

[/quote]
this is a simple paragraph with [i]strong[/i] [b]arguments[/b]!.


[hr][/hr]
Link: [url=https://example.com]foo[/url]
Image: [img]https://fastly.picsum.photos/id/133/200/300.jpg[/img]`

	data := Run(input)
	if string(data) != expected {
		t.Log(string(data))
		t.Errorf("error producing complext output")
	}
}
