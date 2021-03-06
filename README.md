
This is a GO AppEngine Unit Testing (GAET) Framework.
=====================================================

**N O T E:**

**This is very early release and some pretty important functionality is still missing.
The most important thing, in my opinion, is segregating test and development datastore info.
At this point, unit tests will modify your development datastore data, please take this into consideration!!!**

**This should, in theory, be easy to modify once the GO SDK supports datastore namespaces.**

**Also, CompareStructs and ComparePartialStructs need WAY more testing on my part.  I would not, at this point, put much faith in them.**

To use GAET, do the following:

1.  Write your unit tests:<br/>
  All unit tests must contain the signature func x(*gaet.Test). The gaet.Test struct provides an appengine context which can be used in your tests.  If neither gaet.Test.Pass(string) or gaet.Test.Fail(string) are during your test, gaet will assume it was a success. Example Unit Test:

        func TestAddUser(t *gaet.Test) {
          expectedUser := &UserAccount {
            Name:     "Joe User",
            Email:    "joe@testdomain.com",
            UserID:   45404,
            Creation: datastore.SecondsToTime(time.Seconds()),
          }

          joeKey,err := datastore.Put(t.Context, datastore.NewIncompleteKey(t.Context, "UserAccount", nil), expectedUser)
          if err != nil {
            t.Fail(err.String())
            return
          }

          err = datastore.Get(t.Context, joeKey, &actualUser)
          if err != nil {
            t.Fail(err.String())
            return
          }

          // ComparePartialStructs will check only the fields listed in the
          // first paramter.  In this case, we're ignoring the creation date.
          err = gaet.ComparePartialStructs([]string{"Name", "Email", "UserID"}, expectedUser, actualUser)
          if err != nil {
            t.Fail(err.String())
            return
          }
        }

2.  Create a unit testing function endpoint:

        func runTests(w http.ResponseWriter, r *http.Request) {
          gaet.RegisterTest("Example Test", TestAddUser)
          gaet.RunTests(w, r)
        }

3.  Finally, register this endpoint in your init function:

        // I only register this endpoint on development, as I don't want unit tests run in production.
        if appengine.IsDevAppServer() {
          http.HandleFunc("/UnitTest", runTests)
        }

Any questions, comments, concerns, or shortcomings can be emailed to ronniemhowell@gmail.com
