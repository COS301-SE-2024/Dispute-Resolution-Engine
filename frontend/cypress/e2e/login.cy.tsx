describe("Navigation", () => {
  it("should login to the server", () => {
    cy.visit("https://capstone-dre.dns.net.za/login");
    cy.contains("Email").type("sediv39443@alientex.com");
    cy.contains("Password").type("Password1234#");
    cy.get("button").contains("Login").click();
  });
});
