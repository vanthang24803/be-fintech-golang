import { describe, expect, test } from "bun:test";

import { isNavItemActive } from "../../../../src/lib/navigation/matchers";

describe("navigation helpers", () => {
  test("marks an item active when the pathname exactly matches the href", () => {
    expect(isNavItemActive("/dashboard", "/dashboard")).toBe(true);
  });

  test("marks an item active for nested routes under the same section", () => {
    expect(isNavItemActive("/transactions/weekly", "/transactions")).toBe(true);
  });

  test("does not match unrelated prefixes", () => {
    expect(isNavItemActive("/transactional-report", "/transactions")).toBe(false);
  });
});

