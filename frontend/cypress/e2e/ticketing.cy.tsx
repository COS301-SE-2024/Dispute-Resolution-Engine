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
    cy.get('input[placeholder="Search tickets..."]').type('account{enter}');
    /**
     * Pagination
     */
    cy.get('ul').contains('Next').should('be.visible');
    cy.get('button').contains('Next').click();
    cy.get('table').should('be.visible');
    /**
     * Checks Drawer
     */
    let firstSubject: string = ""
    cy.get(':nth-child(1) > :nth-child(1) > a').invoke('text').then((text) => {
      firstSubject = text
      cy.get(':nth-child(1) > :nth-child(1) > a').click()
      cy.get('h2').contains(firstSubject)
      cy.get('textarea').type("New Ticket{enter}")
      cy.get('button').contains('Send').click()
      // cy.get('div').contains("New Ticket")
    })
    // cy.get('button').contains('Filter by').click();
    // cy.get('[role="dialog"]').contains('Filter').should('be.visible');
  });
});
