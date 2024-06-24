describe('Navigation', () => {
	it('should get the disputes', () => {
		cy.visit('http://localhost:3000/login')
		cy.contains('Email').type('alice@smith.co.za')
		cy.contains('Password').type('Password1234#')
		cy.get('button').contains('Login').click()
		cy.visit('http://localhost:3000/disputes')
		cy.get('button').contains('Dispute 0').click()
		cy.contains('Dispute ID')
		cy.contains('Create').click()
	})	
  })