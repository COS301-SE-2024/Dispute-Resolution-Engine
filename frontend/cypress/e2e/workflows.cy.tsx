describe("Navigation", () => {
  it("should login to the server", () => {
    /**
     * Logs in
     */
    cy.visit("/login");
    cy.contains("Email").type("kifane5182@jofuso.com");
    cy.contains("Password").type("Test1234#");
    cy.get("button").contains("Login").click();

    /**
     * Check that the right tings are there
     */
    cy.get(':nth-child(2) > .inline-flex').first().click();
    cy.get('div.flex-col > header.flex > .grow').should('be.visible');
    cy.get('input[placeholder="Search tickets..."]').should('be.visible');
    cy.get('button').contains('Filter by').should('be.visible');
    /**
     * Searching
     */
    cy.get('input[placeholder="Search tickets..."]').type('e{enter}');
    /**
     * Pagination
     */
    cy.get('input[placeholder="Search tickets..."]').clear();
    cy.viewport(1920, 2000)
    cy.get('button').contains('Next').should('be.visible');
    cy.get('table').should('be.visible');
    
  });
});