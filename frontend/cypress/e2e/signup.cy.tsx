describe('Navigation', () => {
	it('should login to the server', () => {
		cy.visit('http://localhost:3000/login')
		cy.contains('Email').type('alice@smith.co.za')
		cy.contains('Password').type('Password1234#')
		cy.contains('Confirm Password').type('Password1234#')
		cy.get('button').contains('Create').click()
		cy.get('p').contains('Something happened')
	})	
  })