export class Version {
  constructor(
    public id:             number,
    public uuid:           string,
    public name:           string,
    public slug:           string,
    public image:          string,
    public replicas:       number,
    public environment_id: number,
  ){}
}
