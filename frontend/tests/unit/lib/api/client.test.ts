import { describe, expect, test } from "bun:test";

import { buildApiUrl, getApiBaseUrl } from "../../../../src/lib/api/client";

describe("API helpers", () => {
  test("falls back to the local Go API base URL when env is missing", () => {
    delete process.env.NEXT_PUBLIC_API_BASE_URL;

    expect(getApiBaseUrl()).toBe("http://localhost:8386");
  });

  test("builds URLs without duplicating slashes", () => {
    process.env.NEXT_PUBLIC_API_BASE_URL = "http://localhost:8386/";

    expect(buildApiUrl("/api/v1/health")).toBe("http://localhost:8386/api/v1/health");
  });

  test("trims surrounding whitespace from the configured base URL", () => {
    process.env.NEXT_PUBLIC_API_BASE_URL = "  http://localhost:8386/api  ";

    expect(getApiBaseUrl()).toBe("http://localhost:8386/api");
  });
});

