import "cypress-file-upload";
describe("Navigation", () => {
  let cookie: string | undefined;
  let time = Date.now();
  let title = "Cypress Test Title" + time;
  let description = "I am creating a cypress test";
  it("should be able to create a dispute", () => {
    cy.visit("https://capstone-dre.dns.net.za/login");
    cy.contains("Email").type("yexiy79682@mvpalace.com");
    cy.contains("Password").type("Password1234#");
    cy.get("button").contains("Login").click();
    cy.wait(2000);
    cy.getCookie("jwt").then((tempCookie) => {
      cookie = tempCookie?.value;
    });
    cy.visit("http://capstone-dre.dns.net.za/disputes");
    cy.get("a").contains("+ Create").click();
    cy.get('input[name="respondentName"]').click().type("Bob Charlie");
    cy.get('input[name="respondentEmail"]').click().type("yexiy79682@mvpalace.com");
    cy.get('input[name="respondentTelephone"]').click().type("0123456789");
    cy.get('input[name="title"]').click().type(title);
    cy.get('textarea[name="summary"').click().type(description);
    cy.fixture("test.txt").then((fileContent) => {
      cy.get('input[type="file"]').attachFile({
        fileContent: fileContent.toString(),
        fileName: "test.txt",
        mimeType: "image/png",
      });
    });
    cy.get(".pt-0 > .inline-flex").click();
    cy.get("span").contains(title);
  });
  it("should be able see the dispute", () => {
    cy.setCookie("jwt", cookie ?? "");
    cy.visit("http://capstone-dre.dns.net.za/disputes");
    cy.get("span").contains(title).click();
    cy.get("h1").contains(title);
    cy.get("p").contains(description);
  });
  it("should not be in the archive", () => {
    cy.setCookie("jwt", cookie ?? "");
    cy.visit("http://capstone-dre.dns.net.za/archive");
    cy.get('input[name="q"]').click().type(title);
    cy.get('button[type="submit"]').click();
  });
});
