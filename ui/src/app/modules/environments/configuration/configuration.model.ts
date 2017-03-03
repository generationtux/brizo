export class ConfigurationModel {
  constructor(
    public id:              number,
    public name:            string,
    public value:           string,
    public environment_id:  number
  ){}
}