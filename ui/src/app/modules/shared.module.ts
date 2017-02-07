import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MastheadComponent } from '../components/masthead/masthead.component';

import { ApplicationService } from './applications/application.service';
import { AuthService } from './auth/auth.service';

@NgModule({
  imports:      [ FormsModule, ReactiveFormsModule ],
  declarations: [ MastheadComponent ],
  exports:      [ MastheadComponent ],
  providers:    [ ApplicationService, AuthService ],
})

export class SharedModule {}
