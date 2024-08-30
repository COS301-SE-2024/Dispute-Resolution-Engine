describe('Archive', () => {
	it('should be accessible', () => {
		cy.visit('https://capstone-dre.dns.net.za/archive')
        cy.get('a').contains('Read More').click()
		cy.contains('Date Filed')
		cy.contains('Decision')
		cy.get('dd').contains('2024')
	})	
  })