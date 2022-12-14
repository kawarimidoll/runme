package edit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditor(t *testing.T) {
	e := new(Editor)
	notebook, err := e.Deserialize(testDataNested)
	require.NoError(t, err)
	result, err := e.Serialize(notebook)
	require.NoError(t, err)
	assert.Equal(
		t,
		string(testDataNestedFlattened),
		string(result),
	)
}

func TestEditor_List(t *testing.T) {
	data := []byte(`1. Item 1
2. Item 2
3. Item 3
`)
	e := new(Editor)
	notebook, err := e.Deserialize(data)
	require.NoError(t, err)

	notebook.Cells[0].Value = "1. Item 1\n2. Item 2\n"

	newData, err := e.Serialize(notebook)
	require.NoError(t, err)
	assert.Equal(
		t,
		`1. Item 1
2. Item 2
`,
		string(newData),
	)

	newData, err = e.Serialize(notebook)
	require.NoError(t, err)
	assert.Equal(
		t,
		`1. Item 1
2. Item 2
`,
		string(newData),
	)
}

func TestEditor_CodeBlock(t *testing.T) {
	t.Run("ProvideGeneratedName", func(t *testing.T) {
		data := []byte("```sh\necho 1\n```\n")
		e := new(Editor)
		notebook, err := e.Deserialize(data)
		require.NoError(t, err)
		cell := notebook.Cells[0]
		assert.Equal(
			t,
			cell.Metadata[prefixAttributeName(internalAttributePrefix, "name")].(string),
			"echo-1",
		)
		// "name" is nil because it was not included in the original snippet.
		assert.Nil(
			t,
			cell.Metadata["name"],
		)
		result, err := e.Serialize(notebook)
		require.NoError(t, err)
		assert.Equal(t, string(data), string(result))
	})

	t.Run("PreserveName", func(t *testing.T) {
		data := []byte("```sh { name=name1 }\necho 1\n```\n")
		e := new(Editor)
		notebook, err := e.Deserialize(data)
		require.NoError(t, err)
		cell := notebook.Cells[0]
		assert.Equal(
			t,
			cell.Metadata[prefixAttributeName(internalAttributePrefix, "name")].(string),
			"name1",
		)
		// "name" is not nil because it was included in the original snippet.
		assert.Equal(
			t,
			cell.Metadata["name"].(string),
			"name1",
		)
		result, err := e.Serialize(notebook)
		require.NoError(t, err)
		assert.Equal(t, string(data), string(result))
	})
}
