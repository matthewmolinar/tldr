- [DONE] B1  Bootstrap Go Fiber Server .... Tests to write  unit – create cmd/api/main_test.go using httptest.NewRequest / app.Test (Fiber helper) to assert 200 response. Guides: Fiber docs


- [DONE] B2 Skeleton POST /api/summarize .... Tests to write unit: handler returns 201 with dummy JSON when body is {"url":"https://example.com"}.


- [DONE] B3 URL Validation helper .... Test to write table-driven unit – loop through cases: valid URL, http:// scheme, oversized HTML (simulate by spinning up local httptest server returning large header). Use Go testing + t.Run. Example of httptest usage: https://golang.cafe/blog/golang-httptest-example.html?utm_source=chatgpt.com


- [DONE] B4 Article Extraction module .... Tests to write: unit – golden-file approach: save a known HTML sample under testdata/article.html; run Extract and assert returned text contains expected phrase & len ≤ 8192 bytes.


- [DONE] B5 LLM client wrapper .... Tests to write: unit – inject http.Client stub that intercepts JSON payload, lets you assert max-prompt size & returns mocked LLM JSON.


- [DOING] B6 Wire Extraction with LLM .... Test to write: integration – spin up Fiber app, stub LLM client (inject via interface) returning canned result, hit /api/summarize, verify pipeline.


- [TODO] Error mapping middleware .... Test to write: unit – create artificial handler returning fiber.ErrUnprocessableEntity, ensure response body matches specification.


