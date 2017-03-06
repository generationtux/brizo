export class Configuration {
  public name:              string|null;
  public value:             string|null;
  public environment_uuid:  string;
  public uuid?:             string|null;

  constructor(props: any = {}){
    Object.assign(this, props);
    this.name  = this.name  || null;
    this.value = this.value || null;
    this.uuid  = this.uuid  || null;
  }
}
