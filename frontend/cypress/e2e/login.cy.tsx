describe("Navigation", () => {
  it("should login to the server", () => {
    cy.visit("/login");
    cy.contains("Email").type(Cypress.env('TEST_USER'));
    cy.contains("Password").type(Cypress.env('TEST_PASSWORD'));
    cy.get("button").contains("Login").click();
  });
});
