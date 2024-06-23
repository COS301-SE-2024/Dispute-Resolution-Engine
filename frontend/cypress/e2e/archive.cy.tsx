describe('Navigation', () => {
	it('should login to the server', () => {
		cy.visit('http://localhost:3000/archive')
        cy.get('form > .flex').type('1')
	})	
  })