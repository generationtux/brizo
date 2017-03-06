import { Application } from '../applications/application.model';

export class Environment {
  constructor(
    public id:             number,
    public name:           string,
    public slug:           string,
    public uuid:           string,
    public application_id: number,
    public application?:   Application,
  ){}
}
