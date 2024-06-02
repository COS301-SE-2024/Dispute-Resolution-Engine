describe('Navigation', () => {
	it('should login to the server', () => {
		cy.visit('http://localhost:3000/disputes')
		cy.get('button').first().click()
		cy.contains('Dispute Details')
		cy.contains("Respondant's Evidence")
	})	
  })