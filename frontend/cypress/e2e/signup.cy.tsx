describe('Navigation', () => {
	it('should signup to the server', () => {
		cy.visit('http://localhost:3000/signup')
		cy.contains('First Name').type('Alice')
		cy.contains('Last Name').type('Smith')
		cy.get('#dateOfBirth')
		cy.contains('Email').type('alice@smith.co.za')
		cy.contains('Password').type('Password1234#')
		cy.contains('Confirm Password').type('Password1234#')
		cy.get('#dateOfBirth').type('1990-06-15')
		cy.get('button').contains('Signup').click()
		// cy.get('p').contains('Something happened')
		
	})	
  })