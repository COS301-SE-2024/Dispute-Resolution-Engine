describe('Archive', () => {
	it('should be accessible', () => {
		cy.visit('https://capstone-dre.dns.net.za/archive')
        cy.get('h1').contains('Resolved Disputes')
	})	
  })