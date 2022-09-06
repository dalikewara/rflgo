package rflgo_test

import (
	"encoding/json"
	"fmt"
	"github.com/dalikewara/rflgo"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type prop struct {
	Prop string
	Name string
}

type log struct {
	Log     string
	Changed string
	Prop    prop
}

type scope struct {
	Name     string
	IsActive int
	Logs     []log
	Prop     *prop
}

type user struct {
	Id     int
	Name   string
	Scopes []*scope
	Prop   prop
}

func generateTestUsers() []*user {
	var lgs []log
	lgs = append(lgs, log{
		Log:     "log",
		Changed: "1 day ago",
		Prop: prop{
			Prop: "prop2",
			Name: "propName1",
		},
	})
	lgs = append(lgs, log{
		Log:     "log",
		Changed: "2 days ago",
		Prop: prop{
			Prop: "prop4",
			Name: "propName3",
		},
	})
	lgs = append(lgs, log{
		Log:     "log",
		Changed: "3 days ago",
		Prop: prop{
			Prop: "prop6",
			Name: "propName5",
		},
	})
	var scps []*scope
	scps = append(scps, &scope{
		Name:     "read",
		IsActive: 1,
		Prop: &prop{
			Prop: "prop1",
			Name: "propName2",
		},
		Logs: lgs,
	})
	scps = append(scps, &scope{
		Name: "create",
		Prop: &prop{
			Prop: "prop3",
			Name: "propName4",
		},
		Logs:     lgs,
		IsActive: 0,
	})
	scps = append(scps, &scope{
		Name:     "update",
		IsActive: 1,
		Prop: &prop{
			Prop: "prop5",
			Name: "propName6",
		},
		Logs: lgs[1:2],
	})
	var usrs []*user
	usrs = append(usrs, &user{
		Id:   1,
		Name: "johndoe",
		Prop: prop{
			Prop: "prop1",
			Name: "propName1",
		},
		Scopes: scps[0:1],
	})
	usrs = append(usrs, &user{
		Id:   2,
		Name: "adamsmith",
		Prop: prop{
			Prop: "prop2",
			Name: "propName2",
		},
		Scopes: scps[0:2],
	})
	usrs = append(usrs, &user{
		Scopes: scps[1:2],
		Name:   "dalikewara",
		Id:     3,
		Prop: prop{
			Prop: "prop3",
			Name: "propName3",
		},
	})
	return usrs
}

func TestCompose_OK(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		type testUser struct {
			Id     int
			Name   string
			Scopes []*scope
			Prop   prop
		}
		s := generateTestUsers()
		var d []*testUser
		err := rflgo.Compose(&d, s)
		assert.Nil(t, err)
		assert.Equal(t, len(d), len(s))
		for i := 0; i < len(d); i++ {
			assert.Equal(t, s[i].Id, d[i].Id)
			assert.Equal(t, s[i].Name, d[i].Name)
			assert.Equal(t, s[i].Prop.Prop, d[i].Prop.Prop)
			assert.Equal(t, s[i].Prop.Name, d[i].Prop.Name)
		}
		sb, _ := json.Marshal(s)
		db, _ := json.Marshal(d)
		assert.Equal(t, string(sb), string(db))
	})

	t.Run("case 2", func(t *testing.T) {
		type testUser struct {
			Id   int
			Name string
		}
		s := generateTestUsers()
		var d []*testUser
		err := rflgo.Compose(&d, s)
		assert.Nil(t, err)
		assert.Equal(t, len(d), len(s))
		for i := 0; i < len(d); i++ {
			assert.Equal(t, s[i].Id, d[i].Id)
			assert.Equal(t, s[i].Name, d[i].Name)
		}
	})

	t.Run("case 3", func(t *testing.T) {
		type userSource struct {
			Id        int
			Name      string
			CreatedAt time.Time
		}
		s := []*userSource{
			{
				Id:        1,
				Name:      "johndoe",
				CreatedAt: time.Now(),
			},
			{
				Id:        2,
				Name:      "adamsmith",
				CreatedAt: time.Now(),
			},
		}
		type userDest struct {
			Name string
		}
		var d []*userDest
		err := rflgo.Compose(&d, s)
		if err != nil {
			panic(err)
		}
		assert.Nil(t, err)
		ex := []*userDest{
			{
				Name: "johndoe",
			},
			{
				Name: "adamsmith",
			},
		}
		assert.Equal(t, len(d), len(ex))
		for i := 0; i < len(d); i++ {
			assert.Equal(t, ex[i].Name, d[i].Name)
		}
		exb, _ := json.Marshal(ex)
		db, _ := json.Marshal(d)
		assert.Equal(t, string(exb), string(db))
	})

	t.Run("case 4", func(t *testing.T) {
		type roleSource struct {
			Permission string
			CreatedAt  time.Time
		}
		type userSource struct {
			Id        int
			Name      string
			Roles     *[]roleSource
			CreatedAt time.Time
		}
		r := &[]roleSource{
			{
				Permission: "create",
				CreatedAt:  time.Now(),
			},
		}
		s := []*userSource{
			{
				Id:        1,
				Name:      "johndoe",
				Roles:     r,
				CreatedAt: time.Now(),
			},
			{
				Id:        2,
				Name:      "dalikewara",
				Roles:     r,
				CreatedAt: time.Now(),
			},
		}
		type roleDest struct {
			Permission string
		}
		type userDest struct {
			Name  string
			Roles *[]roleDest
		}
		var d []*userDest
		err := rflgo.Compose(&d, s)
		if err != nil {
			panic(err)
		}
		assert.Nil(t, err)
		ex := []*userDest{
			{
				Name: "johndoe",
				Roles: &[]roleDest{
					{
						Permission: "create",
					},
				},
			},
			{
				Name: "dalikewara",
				Roles: &[]roleDest{
					{
						Permission: "create",
					},
				},
			},
		}
		assert.Equal(t, len(d), len(ex))
		for i := 0; i < len(d); i++ {
			assert.Equal(t, ex[i].Name, d[i].Name)
		}
		sb, _ := json.Marshal(s)
		exb, _ := json.Marshal(ex)
		db, _ := json.Marshal(d)
		assert.Equal(t, string(exb), string(db))
		fmt.Println(string(sb))
		fmt.Println(string(exb))
		fmt.Println(string(db))
	})

}

func TestSet_ErrValueKindNotMatch(t *testing.T) {
	var d string
	var s *int
	ten := 10
	s = &ten
	err := rflgo.Set(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.NotNil(t, err)
	assert.EqualError(t, err, fmt.Sprintf(rflgo.ErrValueKindNotMatch, "string", "ptr", "string", "*int"))
}

func TestSet_OK(t *testing.T) {
	t.Run("primitive", func(t *testing.T) {
		var d int
		var s int
		s = 10
		err := rflgo.Set(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := 10
		assert.Equal(t, ex, d)
	})

	t.Run("pointer", func(t *testing.T) {
		var d *int
		var s *int
		ten := 10
		s = &ten
		err := rflgo.Set(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := 10
		assert.Equal(t, *&ex, *d)
	})

	t.Run("struct", func(t *testing.T) {
		var d struct {
			Id int
		}
		var s struct {
			Id int
		}
		s = struct {
			Id int
		}{
			Id: 10,
		}
		err := rflgo.Set(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := struct {
			Id int
		}{
			Id: 10,
		}
		assert.Equal(t, ex, d)
	})

	t.Run("slice", func(t *testing.T) {
		var d []struct {
			Id int
		}
		var s []struct {
			Id int
		}
		s = []struct {
			Id int
		}{
			{
				Id: 10,
			},
		}
		err := rflgo.Set(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := []struct {
			Id int
		}{
			{
				Id: 10,
			},
		}
		assert.Equal(t, ex, d)
	})
}

func TestSetSlice_ErrValueKindNotSlice(t *testing.T) {
	var d string
	var s []struct {
		Id int
	}
	s = []struct {
		Id int
	}{
		{
			Id: 10,
		},
	}
	err := rflgo.SetSlice(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.NotNil(t, err)
	assert.EqualError(t, err, fmt.Sprintf(rflgo.ErrValueKindNotSlice, "string", "slice", "string", "[]struct { Id int }"))
}

func TestSetSlice_DestEmpty(t *testing.T) {
	var d []struct {
		Id int
	}
	var s []struct {
		Id int
	}
	s = []struct {
		Id int
	}{
		{
			Id: 10,
		},
	}
	err := rflgo.SetSlice(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.Nil(t, err)
	ex := []struct {
		Id int
	}{
		{
			Id: 10,
		},
	}
	assert.Equal(t, ex, d)
}

func TestSetSlice_DestNotEmpty(t *testing.T) {
	var d []struct {
		Id int
	}
	var s []struct {
		Id int
	}
	nine := []struct {
		Id int
	}{
		{
			Id: 9,
		},
	}
	s = []struct {
		Id int
	}{
		{
			Id: 10,
		},
	}
	d = nine
	assert.Equal(t, nine, d)
	err := rflgo.SetSlice(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.Nil(t, err)
	ex := []struct {
		Id int
	}{
		{
			Id: 10,
		},
	}
	assert.Equal(t, ex, d)
}

func TestSetSlice_SourceEmpty(t *testing.T) {
	var d []struct {
		Id int
	}
	var s []struct {
		Id int
	}
	d = []struct {
		Id int
	}{
		{
			Id: 10,
		},
	}
	err := rflgo.SetSlice(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.Nil(t, err)
	ex := []struct {
		Id int
	}{
		{
			Id: 10,
		},
	}
	assert.Equal(t, ex, d)
}

func TestSetSlice_DestEmptySourceEmpty(t *testing.T) {
	var d []struct {
		Id int
	}
	var s []struct {
		Id int
	}
	err := rflgo.SetSlice(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.Nil(t, err)
	var ex []struct {
		Id int
	}
	assert.Equal(t, ex, d)
}

func TestSetStruct_ErrValueKindNotStruct(t *testing.T) {
	var d string
	var s struct {
		Id int
	}
	s = struct {
		Id int
	}{
		Id: 10,
	}
	err := rflgo.SetStruct(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.NotNil(t, err)
	assert.EqualError(t, err, fmt.Sprintf(rflgo.ErrValueKindNotStruct, "string", "struct", "string", "struct { Id int }"))
}

func TestSetStruct_DestEmpty(t *testing.T) {
	var d struct {
		Id int
	}
	var s struct {
		Id int
	}
	s = struct {
		Id int
	}{
		Id: 10,
	}
	err := rflgo.SetStruct(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.Nil(t, err)
	ex := struct {
		Id int
	}{
		Id: 10,
	}
	assert.Equal(t, ex, d)
}

func TestSetStruct_DestNotEmpty(t *testing.T) {
	var d struct {
		Id int
	}
	var s struct {
		Id int
	}
	nine := struct {
		Id int
	}{
		Id: 9,
	}
	s = struct {
		Id int
	}{
		Id: 10,
	}
	d = nine
	assert.Equal(t, nine, d)
	err := rflgo.SetStruct(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.Nil(t, err)
	ex := struct {
		Id int
	}{
		Id: 10,
	}
	assert.Equal(t, ex, d)
}

func TestSetStruct_SourceEmpty(t *testing.T) {
	var d struct {
		Id int
	}
	var s struct {
		Id int
	}
	d = struct {
		Id int
	}{
		Id: 10,
	}
	err := rflgo.SetStruct(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.Nil(t, err)
	ex := struct {
		Id int
	}{
		Id: 10,
	}
	assert.Equal(t, ex, d)
}

func TestSetStruct_DestEmptySourceEmpty(t *testing.T) {
	var d struct {
		Id int
	}
	var s struct {
		Id int
	}
	err := rflgo.SetStruct(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.Nil(t, err)
	var ex struct {
		Id int
	}
	assert.Equal(t, ex, d)
}

func TestSetPointer_ErrValueKindNotPointer(t *testing.T) {
	var d string
	var s *int
	ten := 10
	s = &ten
	err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
	assert.NotNil(t, err)
	assert.EqualError(t, err, fmt.Sprintf(rflgo.ErrValueKindNotPointer, "string", "ptr", "string", "*int"))
}

func TestSetPointer_DestEmpty(t *testing.T) {
	t.Run("primitive", func(t *testing.T) {
		var d *int
		var s *int
		ten := 10
		s = &ten
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := 10
		assert.Equal(t, *&ex, *d)
	})

	t.Run("struct", func(t *testing.T) {
		var d *struct {
			Id int
		}
		var s *struct {
			Id int
		}
		ten := struct {
			Id int
		}{
			Id: 10,
		}
		s = &ten
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := struct {
			Id int
		}{
			Id: 10,
		}
		assert.Equal(t, *&ex, *d)
	})

	t.Run("slice", func(t *testing.T) {
		var d *[]struct {
			Id int
		}
		var s *[]struct {
			Id int
		}
		ten := []struct {
			Id int
		}{
			{
				Id: 10,
			},
		}
		s = &ten
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := []struct {
			Id int
		}{
			{
				Id: 10,
			},
		}
		assert.Equal(t, *&ex, *d)
	})
}

func TestSetPointer_DestNotEmpty(t *testing.T) {
	t.Run("primitive", func(t *testing.T) {
		var d *string
		var s *string
		nine := "nine"
		ten := "ten"
		s = &ten
		d = &nine
		assert.Equal(t, *&nine, *d)
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := "ten"
		assert.Equal(t, *&ex, *d)
	})

	t.Run("struct", func(t *testing.T) {
		var d *struct {
			Id int
		}
		var s *struct {
			Id int
		}
		nine := struct {
			Id int
		}{
			Id: 9,
		}
		ten := struct {
			Id int
		}{
			Id: 10,
		}
		s = &ten
		d = &nine
		assert.Equal(t, *&nine, *d)
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := struct {
			Id int
		}{
			Id: 10,
		}
		assert.Equal(t, *&ex, *d)
	})

	t.Run("slice", func(t *testing.T) {
		var d *[]struct {
			Id int
		}
		var s *[]struct {
			Id int
		}
		nine := []struct {
			Id int
		}{
			{
				Id: 9,
			},
		}
		ten := []struct {
			Id int
		}{
			{
				Id: 10,
			},
		}
		s = &ten
		d = &nine
		assert.Equal(t, *&nine, *d)
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := []struct {
			Id int
		}{
			{
				Id: 10,
			},
		}
		assert.Equal(t, *&ex, *d)
	})
}

func TestSetPointer_SourceEmpty(t *testing.T) {
	t.Run("primitive", func(t *testing.T) {
		var d *int
		var s *int
		ten := 10
		d = &ten
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := 10
		assert.Equal(t, *&ex, *d)
	})

	t.Run("struct", func(t *testing.T) {
		var d *struct {
			Id int
		}
		var s *struct {
			Id int
		}
		ten := struct {
			Id int
		}{
			Id: 10,
		}
		d = &ten
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := struct {
			Id int
		}{
			Id: 10,
		}
		assert.Equal(t, *&ex, *d)
	})

	t.Run("slice", func(t *testing.T) {
		var d *[]struct {
			Id int
		}
		var s *[]struct {
			Id int
		}
		ten := []struct {
			Id int
		}{
			{
				Id: 10,
			},
		}
		d = &ten
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		ex := []struct {
			Id int
		}{
			{
				Id: 10,
			},
		}
		assert.Equal(t, *&ex, *d)
	})
}

func TestSetPointer_DestEmptySourceEmpty(t *testing.T) {
	t.Run("primitive", func(t *testing.T) {
		var d *int
		var s *int
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		var ex *int
		assert.Equal(t, ex, d)
	})

	t.Run("struct", func(t *testing.T) {
		var d *struct {
			Id int
		}
		var s *struct {
			Id int
		}
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		var ex *struct {
			Id int
		}
		assert.Equal(t, ex, d)
	})

	t.Run("slice", func(t *testing.T) {
		var d *[]struct {
			Id int
		}
		var s *[]struct {
			Id int
		}
		err := rflgo.SetPointer(reflect.ValueOf(&d).Elem(), reflect.ValueOf(s))
		assert.Nil(t, err)
		var ex *[]struct {
			Id int
		}
		assert.Equal(t, ex, d)
	})
}
