import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared.module';

import { AuthGuard } from '../auth/auth.guard';
import { VersionDetailsComponent } from './details/version-details.component';

const versionRoutes: Routes = [
  { path: 'environments/:environment-uuid/versions/:version-uuid', component: VersionDetailsComponent, canActivate: [AuthGuard] },
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forChild(versionRoutes),
    SharedModule,
  ],
  declarations: [
    VersionDetailsComponent,
  ],
  exports: [
    VersionDetailsComponent,
    RouterModule,
  ],
})

export class VersionModule {}
