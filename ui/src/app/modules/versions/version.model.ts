export class Version {
  constructor(
    public id:             number,
    public name:           string,
    public slug:           string,
    public uuid:           string,
    public application_id: number,
    public environment_id: number,
    public image:          string,
    public replicas:       number,
  ){}
}