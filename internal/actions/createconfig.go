/*
   piawgcli
   Copyright (C) 2021  Derek Battams <derek@battams.ca>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package actions

import "gitlab.com/ddb_db/piawgcli/internal/context"

type CreateConfigCmd struct {
	PiaId       string `required help:"PIA user id"`
	PiaPassword string `required help:"PIA password"`
	PiaRegionId string `required help:"PIA region id to connect to; use show-regions command to get the region id"`
}

func (cmd *CreateConfigCmd) Run(ctx *context.Context) error {
	return nil
}
