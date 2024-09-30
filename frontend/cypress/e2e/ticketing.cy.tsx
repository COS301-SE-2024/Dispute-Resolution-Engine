describe("Navigation", () => {
  it("should login to the server", () => {
    /**
     * Logs in
     */
    cy.visit("/admin/login");
    cy.contains("Email").type(Cypress.env('TEST_USER'));
    cy.contains("Password").type(Cypress.env('TEST_PASSWORD'));
    cy.get("button").contains("Login").click();

    /**
     * Check that the right tings are there
     */
    cy.visit("/admin")
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
    /**
     * Filtering
     */
    cy.get('button').contains('Filter by').click()
    cy.get('button').contains('No filter').parent().click()
    cy.get('span').contains('Open').parent().click()
    cy.get('button').contains('Apply').click()

    // cy.get('button').contains('Filter by').click()
    cy.get('button').contains('Open').parent().click()
    cy.get('span').contains('No filter').parent().click()
    cy.get('button').contains('Apply').click()
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
      cy.get('div').contains("New Ticket")
    })
  });
});
