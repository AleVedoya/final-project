export interface IEmail {
  id: string;
  from: string;
  to: string;
  content: string;
  subject: string;
  date: string;
  highlight: string[];
}

export interface IEmailResponse {
    from: string;
    to: string;
    date: string;
    subject: string;
    content: string;
  }
  