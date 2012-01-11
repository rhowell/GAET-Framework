package gaet

import (
    "os"
    "fmt"
    "http"
    "time"
    "reflect"
    "template"
    "appengine"
    "appengine/datastore"
)

type TestResult struct {
    TestName   string
    TestResult string
    TestOutput string
    TestTime   float64 // Milliseconds
}

type TestResultSet struct {
    AppName   string
    Results   []TestResult
    TotalTime float64 // Milliseconds
    Timestamp string
    TestCount int
    TestPass  int
    TestFail  int
}

const NANO_TO_MILLI = 1000000.0

var testList = []TestListEntry{}

func RunTests(w http.ResponseWriter, r *http.Request) {
    ctx := appengine.NewContext(r)
    rs := &TestResultSet{
        AppName: appengine.AppID(ctx),
        Results: []TestResult{},
        TotalTime: 0,
        Timestamp: time.LocalTime().Format(time.RFC1123),
        TestCount: 0,
        TestPass:  0,
        TestFail:  0,
    }

    for _,test := range testList {
        t := new(Test)
        t.Context = ctx
        nanoTime := benchmark(test.Test, t)
        rs.TestCount++

        switch {
        case t.Status == TestState[PASS]:
            rs.TestPass++
            break
        case t.Status == TestState[FAIL]:
            rs.TestFail++
            break
        }
        rs.Results = append(rs.Results, TestResult{
            TestName:   test.Name,
            TestResult: t.Status,
            TestOutput: t.Output,
            TestTime:   nanoTime,
        })
        rs.TotalTime += nanoTime
    }

    testTemplate := template.Must(template.New("test_page").Parse(testResultTemplate))
    testTemplate.Execute(w, rs)
}

func RegisterTest(name string, test func(t *Test)) {
    for _,test := range testList {
        if name == test.Name {
            return
        }
    }
    testList = append(testList, TestListEntry{name, test})
}

func ClearTests() {
    testList = []TestListEntry{}
}

func AssertKeyToEntry(ctx appengine.Context, actual *datastore.Key, expected interface{}) bool {
//    entry := interface{}
    return false
  //  datastore.Get(ctx, actual, )
}

func AssertEntryToEntry(ctx appengine.Context, actual, expected interface{}) bool {
    return false
}

func readFile(filename string) string  {
    f, err := os.Open(filename)
    if err != nil {
        return ""
    }
    defer f.Close()  // f.Close will run when we're finished.

    var result []byte
    buf := make([]byte, 1024)
    for {
    n, err := f.Read(buf[0:])
        result = append(result, buf[0:n]...) // append is discussed later.
        if err != nil {
            if err == os.EOF {
                break
            }
            return ""  // f will be closed if we return here.
        }
    }
    return string(result) // f will be closed if we return here.
}

func benchmark(f func(*Test), t *Test) float64 {
    before := float64(time.Nanoseconds())
    f(t)
    if 0 == len(t.Status) {
        t.Pass(t.Output)
    }
    after := float64(time.Nanoseconds())
    return (after - before)/NANO_TO_MILLI
}


// Grabs the fields out of a struct and prints them out
func TestStuff(src interface{}) string {
    s := reflect.ValueOf(src).Elem()
    result := s.Type().String() + "\n"
    typeOfT := s.Type()

    for i := 0; i < s.NumField(); i++ {
        f := s.Field(i)
        result += fmt.Sprintf("%d: %10s %15v = %#19v Is Direct Comparable: %5v %v\n", i,
            typeOfT.Field(i).Name, f.Type(), f.Interface(), isDirectComparable(f.Type()), f.Type().Kind())
    }
    return result
}

func CompareStructs(src interface{}, dest interface{}) os.Error {
    srcReflect  := reflect.ValueOf(src).Elem()
    destReflect := reflect.ValueOf(dest).Elem()
    srcType  := srcReflect.Type()
    destType := destReflect.Type()

    if srcReflect.Type() != destReflect.Type() {
        return os.NewError(fmt.Sprintf("src is of type %v, whereas dest is of type %v", srcType.String(), destType.String()))
    }

    if srcReflect.NumField() != destReflect.NumField() {
        return os.NewError(fmt.Sprintf("src contains %v fields, whereas dest contains %v fields", srcReflect.NumField(), destReflect.NumField()))
    }

    for iter := 0; iter < srcReflect.NumField(); iter++ {
        srcField  :=  srcReflect.Field(iter)
        destField := destReflect.Field(iter)

        if srcField.Type() != destField.Type() {
            return os.NewError(fmt.Sprintf("src field %v is of type \"%v\", whereas dest field is of type \"%v\"", srcType.Field(iter).Name, srcField.Type().String(), destField.Type().String()))
        }

        if isDirectComparable(srcField.Type()) {
            if srcField.Interface() != destField.Interface() {
                return os.NewError(fmt.Sprintf("src field %v contains value \"%v\", whereas dest field contains value \"%v\"", srcType.Field(iter).Name, srcField.Interface(), destField.Interface()))
            }
        } else {
            // Hak job comparisson:
            srcSlice  := fmt.Sprint(srcField.Interface())
            destSlice := fmt.Sprint(destField.Interface())
            if srcSlice != destSlice {
                return os.NewError(fmt.Sprintf("src field %v contains \"%v\", whereas dest field contains \"%v\"", srcType.Field(iter).Name, srcSlice, destSlice))
            }
        }
    }

    return nil
}

func ComparePartialStructs(fieldNames []string, src interface{}, dest interface{}) os.Error {
    srcReflect  := reflect.ValueOf(src).Elem()
    destReflect := reflect.ValueOf(dest).Elem()

    if srcReflect.Type() != destReflect.Type() {
        return os.NewError(fmt.Sprintf("src is of type %v, whereas dest is of type %v", srcReflect.Type().String(), destReflect.Type().String()))
    }

    if srcReflect.NumField() != destReflect.NumField() {
        return os.NewError(fmt.Sprintf("src contains %v fields, whereas dest contains %v fields", srcReflect.NumField(), destReflect.NumField()))
    }

    for _, field := range fieldNames {
        if 0 == len(field) {
            break
        }
        srcField  := srcReflect.FieldByName(field)
        if 0 == srcField.Kind() {
            return os.NewError(fmt.Sprintf("Could not find field %v in source struct", field))
        }
        destField := destReflect.FieldByName(field)
        if 0 == destField.Kind() {
            return os.NewError(fmt.Sprintf("Could not find field %v in dest struct", field))
        }

        if srcField.Kind() != destField.Kind() {
            return os.NewError(fmt.Sprintf("Field %v in source struct is of type %v, but is of type %v in dest struct", field, srcField.Kind(), destField.Kind()))
        }

        if isDirectComparable(srcField.Type()) {
            if srcField.Interface() != destField.Interface() {
                return os.NewError(fmt.Sprintf("Source field %v contains value %v, whereas dest field contains value %v", field, srcField.Interface(), destField.Interface()))
            }
        } else {
            // Hak job comparisson:
            srcSlice  := fmt.Sprint(srcField.Interface())
            destSlice := fmt.Sprint(destField.Interface())
            if srcSlice != destSlice {
                return os.NewError(fmt.Sprintf("src field %v contains %v, whereas dest field contains %v", field, srcSlice, destSlice))
            }
        }
    }

    return nil
}

// Only checks for types usable by datastore
func isDirectComparable(t reflect.Type) bool {
    kind := t.Kind()

    if reflect.Slice == kind {
        return false
    }
    if reflect.Map == kind {
        return false
    }
    if reflect.Array == kind {
        return false
    }
    return true
}
