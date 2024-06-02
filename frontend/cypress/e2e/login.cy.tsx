describe('Navigation', () => {
	it('should siupgn to the server', () => {
		cy.visit('http://localhost:3000/login')
		cy.contains('Email').type('alice@smith.co.za')
		cy.contains('Password').type('Password1234#')
		cy.get('p').contains('Something happened')
	})	
  })